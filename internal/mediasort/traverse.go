package mediasort

import (
	"context"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dtrejod/goexif/internal/ilog"
	"github.com/dtrejod/goexif/internal/mediatype"
	"go.uber.org/zap"
)

type traverser struct {
	sourceDirectory string

	stopWalkOnError        bool
	allowedFileTypes       []string
	blocklist              []*regexp.Regexp
	useInputMagicSignature bool

	fileHandler    *metadataFileHandler
	extVisitorFunc mediatype.VisitorFunc[map[string]struct{}]
}

// innerTraverseFunc is a function that is called for each file that is traversed.
type innerTraverseFunc = func(logger *zap.Logger, srcMedia mediatype.Format) error

// Run implements Sorter
func (t *traverser) Run(ctx context.Context) error {
	progressTracker := &progressTracker{}
	preScanRunFunc := func(logger *zap.Logger, srcMedia mediatype.Format) error {
		progressTracker.recordMediaVisit(ctx, true)
		return nil
	}
	sortMediaRunFunc := func(logger *zap.Logger, srcMedia mediatype.Format) error {
		progressTracker.recordMediaVisit(ctx, false)
		if err := t.fileHandler.handleMediaFile(ctx, srcMedia); err != nil {
			logger.Warn("Failed to handle file.", zap.Error(err))
			if t.stopWalkOnError {
				return err
			}
		}
		return nil
	}

	ilog.FromContext(ctx).Info("Performing pre-scan for media files...", zap.String("directory", t.sourceDirectory))
	if err := filepath.WalkDir(t.sourceDirectory, t.traverseFunc(ctx, preScanRunFunc)); err != nil {
		return err
	}

	ilog.FromContext(ctx).Info("Sorting media files in directory...", zap.String("directory", t.sourceDirectory))
	if err := filepath.WalkDir(t.sourceDirectory, t.traverseFunc(ctx, sortMediaRunFunc)); err != nil {
		return err
	}

	ilog.FromContext(ctx).Info("Successfully sorted media files.")
	return nil
}

func (t *traverser) traverseFunc(ctx context.Context, innerCallbackFunc innerTraverseFunc) fs.WalkDirFunc {
	return func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		logger := ilog.FromContext(ctx).With(zap.String("path", path))

		if info.IsDir() {
			if t.skipDir(path) {
				logger.Debug("Directory matches blocklist, so skipping entire directory.")
				return fs.SkipDir
			}
			return nil
		}

		if t.skipDir(path) {
			logger.Debug("Path in blocklist, so skipping.")
			return fs.SkipDir
		}

		srcMedia, err := mediatype.NewFormat(path, t.useInputMagicSignature)
		if err != nil {
			logger.Debug("Could not identify file as media file.", zap.Error(err))
			return nil
		}

		visistor := mediatype.FormatWithVisitor[map[string]struct{}](srcMedia)
		aliases, err := visistor.Accept(ctx, t.extVisitorFunc)
		if err != nil {
			logger.Warn("Could not identify file from known media aliases.", zap.Error(err))
			return nil
		}

		logger.Debug("Checking file...")
		if t.skipFile(aliases) {
			logger.Debug("File did not matchfile types allowlist, so skipping.")
			return nil
		}

		return innerCallbackFunc(logger, srcMedia)
	}
}

func (t *traverser) skipDir(path string) bool {
	for _, d := range t.blocklist {
		if d.MatchString(strings.ToLower(path)) {
			return true
		}
	}
	return false
}

func (t *traverser) skipFile(srcMediaAlises map[string]struct{}) bool {
	for _, allowedType := range t.allowedFileTypes {
		_, ok := srcMediaAlises[allowedType]
		if ok {
			return false
		}
	}

	return true
}
