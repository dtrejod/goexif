package cmd

import (
	"os"
	"regexp"

	"github.com/dtrejod/goexif/internal/mediasort"
	"github.com/spf13/cobra"
)

const (
	sourceDirFlagName         = "src-dir"
	destinationDirFlagName    = "dest-dir"
	dryRunFlagName            = "dry-run"
	tsAsFilenameFlagName      = "ts-as-filename"
	modTimeFallbackFlagName   = "fallback-mod-time"
	forceFlagName             = "force"
	stopOnErrorFlagName       = "stop-on-err"
	detectDuplicatesFlagName  = "detect-duplicates"
	magicSignatureInFlagName  = "magic-ext-in"
	magicSignatureOutFlagName = "magic-ext-out"
	fileTypesFlagName         = "file-types"
	blocklistRegexFlagName    = "blocklist-re"
)

var (
	sourceDir         string
	destDir           string
	dryRun            bool
	tsAsFilename      bool
	modTimeFallback   bool
	magicSignatureIn  bool
	detectDuplicates  bool
	force             bool
	magicSignatureOut bool
	stopOnError       bool
	fileTypes         []string
	blocklistRe       []string
)

var sortCmd = &cobra.Command{
	Use:   "sort",
	Short: "Sort mediasort files from their exif/file metadata",
	Run:   sortRun,
}

func sortRun(_ *cobra.Command, _ []string) {
	opts := []mediasort.Option{mediasort.WithSourceDirectory(sourceDir)}
	if destDir != "" {
		opts = append(opts, mediasort.WithDestinationDirectory(destDir))
	}
	if dryRun {
		opts = append(opts, mediasort.WithDryRun())
	}
	if tsAsFilename {
		opts = append(opts, mediasort.WithTimestampAsFilename())
	}
	if modTimeFallback {
		opts = append(opts, mediasort.WithLastModifiedFallback())
	}
	if magicSignatureIn {
		opts = append(opts, mediasort.WithInputFileMagicSignature())
	}
	if magicSignatureOut {
		opts = append(opts, mediasort.WithOutputFileMagicSignature())
	}
	if detectDuplicates {
		opts = append(opts, mediasort.WithDetectDuplicates())
	}
	if force {
		opts = append(opts, mediasort.WithOverwriteExisting())
	}
	if stopOnError {
		opts = append(opts, mediasort.WithStopOnError())
	}
	if len(fileTypes) > 0 {
		opts = append(opts, mediasort.WithFileTypes(fileTypes))
	}
	if len(blocklistRe) > 0 {
		// gracefully handle the no regex case
		if blocklistRe[0] == "" {
			opts = append(opts, mediasort.WithRegexBlocklist([]string{}))
		} else {
			opts = append(opts, mediasort.WithRegexBlocklist(blocklistRe))
		}
	}

	s, err := mediasort.NewSorter(ctx, opts...)
	if err != nil {
		os.Exit(1)
	}

	if err := s.Run(ctx); err != nil {
		os.Exit(1)
	}
}

func init() {
	sortCmd.Flags().StringVarP(&sourceDir, sourceDirFlagName, "s", "", "Source directory to scan for media files")
	sortCmd.Flags().StringVar(&destDir,
		destinationDirFlagName,
		"",
		"Destination directory to move files into. If not specified uses the relative directory where the original file was found")
	sortCmd.Flags().BoolVarP(&dryRun, dryRunFlagName, "n", true, "Do nothing, only show what would happen")
	sortCmd.Flags().BoolVar(&tsAsFilename, tsAsFilenameFlagName, false, "Use timestamp as new filename")
	sortCmd.Flags().BoolVar(&modTimeFallback,
		modTimeFallbackFlagName,
		false,
		"Fallback to using file modified time if no exif data is found")
	sortCmd.Flags().BoolVar(&detectDuplicates,
		detectDuplicatesFlagName,
		false,
		"Gracefully skip moving duplicate image when name conflict in destination directory")
	sortCmd.Flags().BoolVar(&magicSignatureIn,
		magicSignatureInFlagName,
		false,
		"Ignore existing file extension and use magic signature instead when identifying files")
	sortCmd.Flags().BoolVar(&magicSignatureOut,
		magicSignatureOutFlagName,
		false,
		"Ignore existing file extension and use magic signature instead when generating new destination path")
	sortCmd.Flags().BoolVar(&force,
		forceFlagName,
		false,
		"Force overwrite any existing media on naming collision WARN: Use with caution!")
	sortCmd.Flags().BoolVar(&stopOnError, stopOnErrorFlagName, false, "Exit on first error")
	sortCmd.Flags().StringArrayVar(&fileTypes,
		fileTypesFlagName,
		mediasort.DefaultFileTypes,
		"Allowlist of file types to match on. NOTE: When used in conjuction with mag-ext-in, then magic metadata may be used")
	sortCmd.Flags().StringArrayVar(&blocklistRe,
		blocklistRegexFlagName,
		sliceReToString(mediasort.DefaultBlocklist),
		"Regex blocklist that will skip")

	_ = sortCmd.MarkFlagRequired(sourceDirFlagName)
	rootCmd.AddCommand(sortCmd)
}

func sliceReToString(in []*regexp.Regexp) []string {
	out := make([]string, 0, len(in))
	for _, r := range in {
		out = append(out, r.String())
	}
	return out
}
