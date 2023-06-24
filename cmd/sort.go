package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/dtrejod/goexif/internal/media"
	"github.com/spf13/cobra"
)

const (
	sourceDirFlagName         = "source-dir"
	destinationDirFlagName    = "dest-dir"
	dryRunFlagName            = "dry-run"
	tsAsFilenameFlagName      = "ts-as-filename"
	modTimeFallbackFlagName   = "fallback-mod-time"
	overwriteExistingFlagName = "overwrite"
	stopOnErrorFlagName       = "stop-on-err"
	magicSignatureInFlagName  = "magic-ext-in"
	magicSignatureOutFlagName = "magic-ext-out"
	fileExtsFlagName          = "extensions"
	blocklistRegexFlagName    = "blocklist-regexes"
)

var (
	sourceDir         string
	destDir           string
	dryRun            bool
	tsAsFilename      bool
	modTimeFallback   bool
	magicSignatureIn  bool
	overwrite         bool
	magicSignatureOut bool
	stopOnError       bool
	fileExts          []string
	blocklistRe       []string
)

var sortCmd = &cobra.Command{
	Use:   "sort",
	Short: "Sort media files from their exif/file metadata",
	Run:   shortRun,
}

func shortRun(_ *cobra.Command, _ []string) {
	opts := []media.Option{media.WithSourceDirectory(sourceDir)}
	if destDir != "" {
		opts = append(opts, media.WithDestinationDirectory(destDir))
	}
	if dryRun {
		opts = append(opts, media.WithDryRun())
	}
	if tsAsFilename {
		opts = append(opts, media.WithTimestampAsFilename())
	}
	if modTimeFallback {
		opts = append(opts, media.WithLastModifiedFallback())
	}
	if magicSignatureIn {
		opts = append(opts, media.WithIdentifyFileMagicSignature())
	}
	if magicSignatureOut {
		opts = append(opts, media.WithGenOutputFileMagicSignature())
	}
	if overwrite {
		opts = append(opts, media.WithOverwriteExisting())
	}
	if stopOnError {
		opts = append(opts, media.WithStopOnError())
	}
	if len(fileExts) > 0 {
		opts = append(opts, media.WithFileTypes(fileExts))
	}
	if len(blocklistRe) > 0 {
		opts = append(opts, media.WithRegexBlocklist(blocklistRe))
	}

	s, err := media.NewSorter(ctx, opts...)
	if err != nil {
		os.Exit(1)
	}

	if err := s.Run(ctx); err != nil {
		os.Exit(1)
	}

	return
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
	sortCmd.Flags().BoolVar(&magicSignatureIn,
		magicSignatureInFlagName,
		false,
		"Ignore existing file extension and use magic signature instead when identifying files")
	sortCmd.Flags().BoolVar(&magicSignatureOut,
		magicSignatureOutFlagName,
		false,
		"Ignore existing file extension and use magic signature instead when generating new destination path")
	sortCmd.Flags().BoolVar(&overwrite,
		overwriteExistingFlagName,
		false,
		"Overwrite existing files on rename. WARN: Use with caution!")
	sortCmd.Flags().BoolVar(&stopOnError, stopOnErrorFlagName, false, "Exit on first error")
	sortCmd.Flags().StringArrayVar(&fileExts,
		fileExtsFlagName,
		media.DefaultFileTypes,
		"Allowlist of file extensions to match on")
	sortCmd.Flags().StringArrayVar(&blocklistRe,
		blocklistRegexFlagName,
		sliceReToString(media.DefaultBlocklist),
		fmt.Sprintf("Regex blocklist that will skip"))

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
