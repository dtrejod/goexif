package media

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dtrejod/goexif/internal/exifdata"
	"github.com/dtrejod/goexif/internal/ilog"
	"github.com/h2non/filetype"
	"go.uber.org/zap"
)

const (
	outPathDateFormat = "2006/01/02"
)

var (
	errUnknownMedia = errors.New("unhandled media type")
	errRuntime      = errors.New("runtime error")
)

type sorter struct {
	dryRun              bool
	timestampAsFilename bool
	useLastModifiedDate bool
	useMagicSignature   bool
	cleanFileExtensions bool
	stopWalkOnError     bool
	overwriteExisting   bool

	fileTypes            []string
	blocklist            []*regexp.Regexp
	sourceDirectory      string
	destinationDirectory *string
}

// Run implements Sorter
func (s *sorter) Run(ctx context.Context) error {
	ilog.FromContext(ctx).Info("Sorting media files in directory.", zap.String("directory", s.sourceDirectory))
	if err := filepath.WalkDir(s.sourceDirectory, s.traverseFunc(ctx)); err != nil {
		return err
	}

	ilog.FromContext(ctx).Info("Succesfully sorted media files.")
	return nil
}

func (s *sorter) traverseFunc(ctx context.Context) fs.WalkDirFunc {
	return func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		logger := ilog.FromContext(ctx).With(zap.String("path", path))

		if info.IsDir() {
			if s.isBlocked(path) {
				logger.Debug("Directory matches blocklist, so skipping entire directory...")
				return fs.SkipDir
			}
			return nil
		}

		if s.isBlocked(path) {
			logger.Debug("Path in blocklist, so skipping...")
			return fs.SkipDir
		}

		logger.Debug("Checking file.")
		if !s.isFileTypeMatch(ctx, path) {
			logger.Debug("File did not match handled file types, so skipping...")
			return nil
		}

		if err := s.handleFile(ctx, path); err != nil {
			logger.Warn("Failed to handle file.", zap.Error(err))
			if s.stopWalkOnError {
				return err
			}
		}

		return nil
	}
}

func (s *sorter) isBlocked(path string) bool {
	for _, d := range s.blocklist {
		if d.MatchString(strings.ToLower(path)) {
			return true
		}
	}
	return false
}

func (s *sorter) isFileTypeMatch(ctx context.Context, path string) bool {
	ext, err := getExt(path, s.useMagicSignature)
	if err != nil {
		return false
	}
	ext = strings.TrimPrefix(ext, ".")
	for _, t := range s.fileTypes {
		if ext == t {
			return true
		}
	}

	ilog.FromContext(ctx).Debug("Unhandled extension", zap.String("ext", ext))
	return false
}

func (s *sorter) handleFile(ctx context.Context, srcPath string) error {
	logger := ilog.FromContext(ctx).With(zap.String("sourcePath", srcPath))

	ts, err := s.getFileTimestamp(srcPath)
	if err != nil {
		return err
	}
	logger.Debug("Discovered timestamp.", zap.String("timestamp", ts.String()), zap.Time("timestampUnix", ts))

	outPath, err := s.getOutputFile(ts, srcPath)
	if err != nil {
		return err
	}
	logger = logger.With(zap.String("destinationPath", outPath))

	if srcPath == outPath {
		logger.Info("Source and destination file match. Nothing to do.")
		return nil
	}

	if s.dryRun {
		logger.Info("Dry run, moving file...")
		return nil
	}
	logger.Debug("Moving file...")

	if err := s.handleOverwrite(ctx, outPath); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}

	if err := os.Rename(srcPath, outPath); err != nil {
		return err
	}

	logger.Info("Successfully moved file.")
	return nil
}

func (s *sorter) getFileTimestamp(path string) (ts time.Time, err error) {
	defer func() {
		// on error, fallback to lastmodified if the option was specified
		if err != nil && s.useLastModifiedDate {
			f, statErr := os.Stat(path)
			if statErr != nil {
				return
			}
			ts = f.ModTime()
			err = nil
		}
	}()

	t, err := filetype.MatchFile(path)
	if err != nil {
		return time.Time{}, err
	}
	// Only use exif metadata if the file is an image
	if v := t.MIME.Type; v != "image" && v != "bitmap" {
		return time.Time{}, fmt.Errorf("%w: %s", errUnknownMedia, v)
	}

	ts, err = exifdata.GetExifTime(path)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %w", errRuntime, err)
	}
	return ts.UTC(), nil
}

func (s *sorter) getOutputFile(ts time.Time, srcPath string) (string, error) {
	origDir := filepath.Dir(srcPath)
	outDir := filepath.Join(origDir, ts.Format(outPathDateFormat))
	if s.destinationDirectory != nil {
		outDir = filepath.Join(*s.destinationDirectory, ts.Format(outPathDateFormat))
	}

	ext, err := getExt(srcPath, s.cleanFileExtensions)
	if err != nil {
		return "", fmt.Errorf("%w: could not get file extension", err)
	}
	outFilename := strings.TrimSuffix(filepath.Base(srcPath), filepath.Ext(srcPath))
	if s.timestampAsFilename {
		outFilename = strconv.FormatInt(ts.Unix(), 10)
	}
	if outFilename == "" {
		return "", fmt.Errorf("%w: output file has no filename", errRuntime)
	}
	outFilename = outFilename + ext

	return filepath.Join(outDir, outFilename), nil
}

func (s *sorter) handleOverwrite(ctx context.Context, path string) error {
	_, err := os.Stat(path)
	if err == nil {
		// file exists
		if s.overwriteExisting {
			ilog.FromContext(ctx).Info("Overwrite flag set. Removing existing file.",
				zap.String("path", path))
			// delete existing file it exists
			return os.Remove(path)
		}
		return fmt.Errorf("%w: desired destination file already exists", errRuntime)
	}
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	// unknown error
	return fmt.Errorf("%w: could not check if desired out file exist", err)
}

// getExt returns the file extension, using the magic signature if desired. The
// returned file extension includes the period suffix.
func getExt(path string, useSignature bool) (string, error) {
	ext := filepath.Ext(strings.ToLower(path))
	if useSignature {
		t, err := filetype.MatchFile(path)
		if err != nil {
			return "", err
		}
		ext = "." + t.Extension
	}
	if ext == "" {
		return "", errUnknownMedia
	}
	return ext, nil
}
