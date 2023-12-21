package media

import (
	"context"
	"errors"
	"fmt"
	"image"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/corona10/goimagehash"
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
	detectDuplicates    bool
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

	skip, err := s.shouldSkip(ctx, srcPath, outPath)
	if err != nil {
		return err
	}
	if skip {
		ilog.FromContext(ctx).Info("Skipping file...")
		return nil
	}

	if s.dryRun {
		logger.Info("Dry run, moving file...")
		return nil
	}

	logger.Debug("Moving file...")
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

func (s *sorter) shouldSkip(ctx context.Context, srcPath, outPath string) (bool, error) {
	_, err := os.Stat(outPath)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("%w: could not check if desired out file exist", err)
	}

	// file exists. Handle gracefully...
	if s.detectDuplicates {
		isDuplicate, err := s.isDuplicateImage(ctx, srcPath, outPath)
		if err != nil {
			return false, fmt.Errorf("%w: failed to detect if images are duplicate", err)
		}
		if isDuplicate {
			ilog.FromContext(ctx).Info("Detected duplicate image.")
			return true, nil
		}
	}

	if s.overwriteExisting && !s.dryRun {
		ilog.FromContext(ctx).Info("Overwrite flag set. Removing existing file.",
			zap.String("outPath", outPath))
		return false, os.Remove(outPath)
	}
	return false, fmt.Errorf("%w: desired output filename collision", errRuntime)
}

func (s *sorter) isDuplicateImage(ctx context.Context, srcPath, outPath string) (bool, error) {
	logger := ilog.FromContext(ctx).With(
		zap.String("sourcePath", srcPath),
		zap.String("destinationPath", outPath))

	logger.Info("Detecting if images are duplicates...")
	hashA, err := getImagePerceptionHash(srcPath)
	if err != nil {
		return false, err
	}

	hashB, err := getImagePerceptionHash(outPath)
	if err != nil {
		return false, err
	}

	distance, err := hashA.Distance(hashB)
	if err != nil {
		return false, fmt.Errorf("%w: failed to compare images for duplicates", err)
	}

	logger.Debug("Comparing images...",
		zap.String("hashA", hashA.ToString()),
		zap.String("hashB", hashB.ToString()),
		zap.Int("distance", distance))

	// A hash distance > 50 is the threshold for different images.
	return distance < 25, nil
}

// getImagePerceptionHash returns the peception hash of the image. The phash is
// taken first performing a discrete cosine transform on the image. Then it
// compares each pixel to it's average. If it is larger then output 1 else 0
// otherwise.
// Since a Phash operates in the frequency domain, it should be more tolerant
// to color shifts, size scaling, watermarks, and compression between 2 images.
func getImagePerceptionHash(path string) (*goimagehash.ImageHash, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to open image", err)
	}

	// Noteably image.Decode does not work with HEIF encoded images
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to decode image", err)
	}

	hashA, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get perception hash image", err)
	}
	return hashA, nil
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
