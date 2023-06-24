package mediasorter

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
	unknownMediaErr = errors.New("unhandled media type")
)

type sorter struct {
	dryRun              bool
	timestampAsFilename bool
	useLastModifiedDate bool
	useMagicSignature   bool
	cleanFileExtensions bool
	stopWalkOnError     bool

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
				logger.Debug("Directory matches blocklist. Skipping entire directory...")
				return fs.SkipDir
			}
			return nil
		}

		if s.isBlocked(path) {
			logger.Debug("Path in blocklist, so skipping...")
			return fs.SkipDir
		}

		logger.Debug("Checking file.")
		if !s.isExtMatch(path) {
			logger.Debug("File didn't match handled file types.")
			return nil
		}

		if err := s.handleFile(ctx, path); err != nil {
			logger.Error("Failed to handle file", zap.Error(err))
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

func (s *sorter) isExtMatch(path string) bool {
	ext, err := getExt(path, s.useMagicSignature)
	if err != nil {
		return false
	}
	for _, t := range s.fileTypes {
		if ext == t {
			return true
		}
	}

	return false
}

func (s *sorter) handleFile(ctx context.Context, srcPath string) error {
	logger := ilog.FromContext(ctx).With(zap.String("path", srcPath))
	logger.Info("Checking EXIF for file.")

	ts, err := s.getFileTimestamp(srcPath)
	if err != nil {
		return err
	}
	logger.Info("Discovered timestamp", zap.Time("timestamp", ts))

	outFileName, err := s.getOutputFile(ts, srcPath)
	if err != nil {
		return err
	}

	if s.dryRun {
		logger.Info("Dry run... Mock moving file.",
			zap.String("sourcePath", srcPath),
			zap.String("destinationPath", outFileName))
	}
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
		return time.Time{}, fmt.Errorf("%w: %s", unknownMediaErr, v)
	}

	return exifdata.GetExifTime(path)
}

func (s *sorter) getOutputFile(ts time.Time, srcPath string) (string, error) {
	origDir := filepath.Dir(srcPath)
	outDir := filepath.Join(origDir, ts.Format(outPathDateFormat))
	if s.destinationDirectory != nil {
		outDir = filepath.Join(*s.destinationDirectory, outDir)
	}

	ext, err := getExt(srcPath, s.cleanFileExtensions)
	if err != nil {
		return "", fmt.Errorf("%w: could not get file extension", err)
	}
	outFilename := filepath.Base(srcPath)
	if s.timestampAsFilename {
		outFilename = strconv.FormatInt(ts.Unix(), 10)
	}
	outFilename = outFilename + ext

	return filepath.Join(outDir, outFilename), nil
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
		return "", unknownMediaErr
	}
	return ext, nil
}
