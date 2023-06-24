package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/dtrejod/goexif/internal/mediasorter"
	"github.com/spf13/cobra"
)

const (
	sourceDirFlagName         = "source-dir"
	destinationDirFlagName    = "dest-dir"
	dryRunFlagName            = "dry-run"
	tsAsFilenameFlagName      = "ts-as-filename"
	modTimeFallbackFlagName   = "fallback-mod-time"
	magicSignatureFlagName    = "magic-sig"
	overwriteExistingFlagName = "overwrite"
	stopOnErrorFlagName       = "stop-on-err"
	cleanFileExtFlagName      = "clean-file-ext"
	fileExtsFlagName          = "extensions"
	blocklistRegexFlagName    = "blocklist-regexes"
)

var (
	sourceDir       string
	destDir         string
	dryRun          bool
	tsAsFilename    bool
	modTimeFallback bool
	magicSignature  bool
	overwrite       bool
	cleanFileExt    bool
	stopOnError     bool
	fileExts        []string
	blocklistRe     []string
)

var sortCmd = &cobra.Command{
	Use:   "sort",
	Short: "Sort media files from their exif/file metadata",
	Run:   shortRun,
}

func shortRun(_ *cobra.Command, _ []string) {
	opts := []mediasorter.Option{mediasorter.WithSourceDirectory(sourceDir)}
	if destDir != "" {
		opts = append(opts, mediasorter.WithDestinationDirectory(destDir))
	}
	if dryRun {
		opts = append(opts, mediasorter.WithDryRun())
	}
	if tsAsFilename {
		opts = append(opts, mediasorter.WithTimestampAsFilename())
	}
	if modTimeFallback {
		opts = append(opts, mediasorter.WithLastModifiedFallback())
	}
	if magicSignature {
		opts = append(opts, mediasorter.WithUseFileMagicSignature())
	}
	if cleanFileExt {
		opts = append(opts, mediasorter.WithCleanFileExtensions())
	}
	if overwrite {
		opts = append(opts, mediasorter.WithOverwriteExisting())
	}
	if stopOnError {
		opts = append(opts, mediasorter.WithStopOnError())
	}
	if len(fileExts) > 0 {
		opts = append(opts, mediasorter.WithFileTypes(fileExts))
	}
	if len(blocklistRe) > 0 {
		opts = append(opts, mediasorter.WithRegexBlocklist(blocklistRe))
	}

	s, err := mediasorter.NewSorter(ctx, opts...)
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
	sortCmd.Flags().BoolVar(&magicSignature,
		magicSignatureFlagName,
		false,
		"Ignore existing file extension and use magic signature instead when identifying files")
	sortCmd.Flags().BoolVar(&cleanFileExt,
		cleanFileExtFlagName,
		false,
		"Ignore existing file extension and use magic signature instead when generating new destination path")
	sortCmd.Flags().BoolVar(&overwrite,
		overwriteExistingFlagName,
		false,
		"Overwrite existing files on rename. WARN: Use with caution!")
	sortCmd.Flags().BoolVar(&stopOnError, stopOnErrorFlagName, false, "Exit on first error")
	sortCmd.Flags().StringArrayVar(&fileExts,
		fileExtsFlagName,
		mediasorter.DefaultFileTypes,
		"Allowlist of file extensions to match on")
	sortCmd.Flags().StringArrayVar(&blocklistRe,
		blocklistRegexFlagName,
		sliceReToString(mediasorter.DefaultBlocklist),
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
