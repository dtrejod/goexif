package mediasort

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/dtrejod/goexif/internal/ilog"
	"github.com/dtrejod/goexif/internal/visitors"
	"go.uber.org/zap"
)

var (
	// DefaultFileTypes are the default media types handled by the sorter if none are specified.
	// NOTE: This default list should match the known mediatypes in the ./internal/mediatype package.
	DefaultFileTypes = []string{
		"jpg",
		"jpeg",
		"png",
		"heif",
	}
	// DefaultBlocklist are default regexes that are ignored by the sorter.
	DefaultBlocklist = []*regexp.Regexp{
		// ignore date pattern used internally by the sorter.
		// Makes subsequent runs iterative instead of
		// repeatedly processing the same media files
		regexp.MustCompile(`(\/)?\d{4}\/(\d{2}\/){2}`),
	}
	errInvalidConfig = errors.New("invalid configuration")
)

// Sorter sorts media from file metadata
type Sorter interface {
	// Run runs the sorter
	Run(ctx context.Context) error
}

// Option is a param that can be used to configure the media metadata sorter.
type Option interface {
	apply(*builderOptions) error
}

type builderFunc func(*builderOptions) error

func (f builderFunc) apply(b *builderOptions) error {
	return f(b)
}

type builderOptions struct {
	dryRun              bool
	timestampAsFilename bool
	useLastModifiedDate bool
	useMagicSignature   bool
	cleanFileExtensions bool
	overwriteExisting   bool
	stopWalkOnError     bool
	detectDuplicates    bool

	allowedFileTypes []string
	blocklist        []*regexp.Regexp

	sourceDirectory      *string
	destinationDirectory *string
}

// NewSorter returns a sorter configured with the provided Option(s). The
// WithSourceDirectory Option is the only required option.
func NewSorter(ctx context.Context, opts ...Option) (Sorter, error) {
	cfg := builderOptions{
		allowedFileTypes: uniqLoweredSlice(DefaultFileTypes),
		blocklist:        DefaultBlocklist,
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt.apply(&cfg); err != nil {
			return nil, err
		}
	}

	if cfg.sourceDirectory == nil {
		err := fmt.Errorf("%w: source directory required for sorting", errInvalidConfig)
		ilog.FromContext(ctx).Error("Failed to build sorter", zap.Error(err))
		return nil, err
	}

	ilog.FromContext(ctx).Info("Sorter configuration.", zap.String("configuration", fmt.Sprintf("%+v", cfg)))
	return &traverser{
		useMagicSignature: cfg.useMagicSignature,
		stopWalkOnError:   cfg.stopWalkOnError,
		allowedFileTypes:  cfg.allowedFileTypes,
		blocklist:         cfg.blocklist,
		sourceDirectory:   *cfg.sourceDirectory,

		extVisitorFunc: visitors.NewMediaExtAliases(ctx),
		fileHandler: &metadataFileHandler{
			useMagicSignature: cfg.useMagicSignature,
			detectDuplicates:  cfg.detectDuplicates,
			dryRun:            cfg.dryRun,
			overwriteExisting: cfg.overwriteExisting,
			mediaMetadataVisitorFunc: visitors.NewMediaMetadataFilename(
				ctx,
				cfg.destinationDirectory,
				cfg.useLastModifiedDate,
				cfg.timestampAsFilename,
				cfg.cleanFileExtensions,
			),
		},
	}, nil
}

// WithDryRun instructs the sorter to make no changes
func WithDryRun() Option {
	return builderFunc(func(b *builderOptions) error {
		b.dryRun = true
		return nil
	})
}

// WithTimestampAsFilename instructs the sorter to rename the source file using it's timestamp and file extension.
// Note: This option can help eliminate duplicate images during sorting.
func WithTimestampAsFilename() Option {
	return builderFunc(func(b *builderOptions) error {
		b.timestampAsFilename = true
		return nil
	})
}

// WithLastModifiedFallback instructs the sorter to fallback to using the
// file's last modified date if there is no media metadata. If false, images without
// media metadata data are ignored
func WithLastModifiedFallback() Option {
	return builderFunc(func(b *builderOptions) error {
		b.useLastModifiedDate = true
		return nil
	})
}

// WithIdentifyFileMagicSignature instructs the sorter to idenitify media files using the
// file's magic signature. If set, then files are renamed to the approriate
// extension accordingly.
// See the manual page for file(1) to understand how this works.
func WithIdentifyFileMagicSignature() Option {
	return builderFunc(func(b *builderOptions) error {
		b.useMagicSignature = true
		return nil
	})
}

// WithGenOutputFileMagicSignature will cause media file extensions to be consistent.
// For example, .jpeg will be renamed to .jpg
func WithGenOutputFileMagicSignature() Option {
	return builderFunc(func(b *builderOptions) error {
		b.cleanFileExtensions = true
		return nil
	})
}

// WithFileTypes is an array of filetypes that we intend to locate.
// Extensions are matched case-insensitive. *.jpg is treated the same as *.JPG, etc.
// Can handle any file type; not just EXIF-enabled file types when used in conjunction with WithUseLastModifiedDate().
func WithFileTypes(t []string) Option {
	return builderFunc(func(b *builderOptions) error {
		b.allowedFileTypes = uniqLoweredSlice(t)
		return nil
	})
}

// WithRegexBlocklist is an array of regular expressions for matching on
// paths to ignore when finding folders. Directory are matched
// case-insensitive
func WithRegexBlocklist(d []string) Option {
	return builderFunc(func(b *builderOptions) error {
		patterns := uniqLoweredSlice(d)

		exprs := make([]*regexp.Regexp, 0, len(patterns))
		for _, p := range patterns {
			re, err := regexp.Compile(p)
			if err != nil {
				return err
			}
			exprs = append(exprs, re)
		}
		b.blocklist = exprs
		return nil
	})
}

// WithSourceDirectory is an absolute or relative filepath where sorted media will looked for.
func WithSourceDirectory(s string) Option {
	return builderFunc(func(b *builderOptions) error {
		b.sourceDirectory = &s
		return nil
	})
}

// WithDestinationDirectory is an absolute or relative filepath where sorted
// media will be saved to.
func WithDestinationDirectory(d string) Option {
	return builderFunc(func(b *builderOptions) error {
		b.destinationDirectory = &d
		return nil
	})
}

// WithOverwriteExisting instructs the sorter to overwrite any existing files
// that may already exist with the same desired destination file name
// Warning: Can be useful for removing duplicates by ensuring no two files with
// the same timestamp can exist, however, can cause data loss if not careful
func WithOverwriteExisting() Option {
	return builderFunc(func(b *builderOptions) error {
		b.overwriteExisting = true
		return nil
	})
}

// WithStopOnError instructs the sorter to exit quickly when any error occurs during walking the directory tree
func WithStopOnError() Option {
	return builderFunc(func(b *builderOptions) error {
		b.stopWalkOnError = true
		return nil
	})
}

// WithDetectDuplicates will use perception hash algorithm of each file to
// determine whether to images with the same EXIF metadata are duplicate files.
func WithDetectDuplicates() Option {
	return builderFunc(func(b *builderOptions) error {
		b.detectDuplicates = true
		return nil
	})
}

// uniqLoweredSlice takes a slice, lowercases all elements, and return a resulting slice with only unique elements.
func uniqLoweredSlice(in []string) []string {
	m := make(map[string]struct{}, len(in))
	for _, s := range in {
		m[strings.ToLower(s)] = struct{}{}
	}

	out := make([]string, 0, len(m))
	for s := range m {
		out = append(out, s)
	}
	return out
}