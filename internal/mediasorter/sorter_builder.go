package mediasorter

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/dtrejod/goexif/internal/ilog"
	"go.uber.org/zap"
)

var (
	// DefaultFileTypes are the default media types handled by the sorter if none are specified.
	DefaultFileTypes = []string{
		".jpg",
		".jpeg",
		".png",
		".tif",
		".tiff",
		".gif",
		".xcf",
	}
	// DefaultBlocklist
	DefaultBlocklist = []*regexp.Regexp{
		regexp.MustCompile(`(\/)?lost\/`),
		regexp.MustCompile(`(\/)?noexif\/`),
		regexp.MustCompile(`(\/)?duplicates\/`),
		regexp.MustCompile(`(\/)?slideshows\/`),
		regexp.MustCompile(`(\/)?raw\/`),
	}
	invalidConfigErr = errors.New("invalid configuration")
)

type Sorter interface {
	Run(ctx context.Context) error
}

// Option is a param that can be used to configure the exif sorter.
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
	stopWalkOnError     bool

	fileTypes []string
	blocklist []*regexp.Regexp

	sourceDirectory      *string
	destinationDirectory *string
}

func NewSorter(ctx context.Context, opts ...Option) (Sorter, error) {
	cfg := &builderOptions{
		fileTypes: uniqLoweredSlice(DefaultFileTypes),
		blocklist: DefaultBlocklist,
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt.apply(cfg); err != nil {
			return nil, err
		}
	}

	if cfg.sourceDirectory == nil {
		err := fmt.Errorf("%w: source directory required for sorting", invalidConfigErr)
		ilog.FromContext(ctx).Error("Failed to build sorter", zap.Error(err))
		return nil, err
	}

	ilog.FromContext(ctx).Info("Sorter configuration.", zap.Reflect("configuration", cfg))
	return &sorter{
		dryRun:               cfg.dryRun,
		timestampAsFilename:  cfg.timestampAsFilename,
		useLastModifiedDate:  cfg.useLastModifiedDate,
		useMagicSignature:    cfg.useMagicSignature,
		cleanFileExtensions:  cfg.cleanFileExtensions,
		stopWalkOnError:      cfg.stopWalkOnError,
		fileTypes:            cfg.fileTypes,
		blocklist:            cfg.blocklist,
		sourceDirectory:      *cfg.sourceDirectory,
		destinationDirectory: cfg.destinationDirectory,
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
// file's last modified date if there is no exif data. If false, images without
// exif data are ignored
func WithLastModifiedFallback() Option {
	return builderFunc(func(b *builderOptions) error {
		b.useLastModifiedDate = true
		return nil
	})
}

// WithUseFileMagicSignature instructs the sorter to idenitify media files using the
// file's magic signature. If set, then files are renamed to the approriate
// extension accordingly.
// See the manual page for file(1) to understand how this works.
func WithUseFileMagicSignature() Option {
	return builderFunc(func(b *builderOptions) error {
		b.useMagicSignature = true
		return nil
	})
}

// WithCleanFileExtensions will cause media file extensions to be consistent.
// For example, .jpeg will be renamed to .jpg
func WithCleanFileExtensions() Option {
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
		b.fileTypes = uniqLoweredSlice(t)
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

// WithStopOnError instructs the sorter to exit quickly when any error occurs during walking the directory tree
func WithStopOnError() Option {
	return builderFunc(func(b *builderOptions) error {
		b.stopWalkOnError = true
		return nil
	})
}

// uniqLoweredSlice takes a slice, lowecases all elements, and return a resulting slice with only unique elements.
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
