package cmd

import (
	"fmt"
	"os"

	"github.com/dtrejod/goexif/internal/mediasorter"
	"github.com/spf13/cobra"
)

const (
	sourceDirFlagName       = "source-directory"
	destinationDirFlagName  = "dest-directory"
	dryRunFlagName          = "dry-run"
	tsAsFilenameFlagName    = "ts-as-filename"
	modTimeFallbackFlagName = "mod-time-fallback"
	magicSignatureFlagName  = "magic-sig"
	stopOnErrorFlagName     = "quick-exit"
	cleanFileExtFlagName    = "clean-ext"
	fileExtsFlagName        = "extensions"
)

var (
	sourceDir       string
	destDir         string
	dryRun          bool
	tsAsFilename    bool
	modTimeFallback bool
	magicSignature  bool
	cleanFileExt    bool
	stopOnError     bool
	fileExts        []string
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
	if stopOnError {
		opts = append(opts, mediasorter.WithStopOnError())
	}
	if len(fileExts) > 0 {
		opts = append(opts, mediasorter.WithFileTypes(fileExts))
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
	sortCmd.Flags().StringVarP(&sourceDir, sourceDirFlagName, "s", "", "source directory to scan for media files")
	sortCmd.Flags().StringVar(&destDir, destinationDirFlagName, "", "destination directory to move files into")
	sortCmd.Flags().BoolVarP(&dryRun, dryRunFlagName, "n", true, "Do nothing, only show what would happen")
	sortCmd.Flags().BoolVar(&tsAsFilename, tsAsFilenameFlagName, false, "Use timestamp as new filename")
	sortCmd.Flags().BoolVar(&modTimeFallback,
		modTimeFallbackFlagName,
		false,
		"Fallback to using file modified time if no exif data is found")
	sortCmd.Flags().BoolVar(&magicSignature,
		magicSignatureFlagName,
		false,
		"Ignore existing file extension and use magic signature instead")
	sortCmd.Flags().BoolVar(&cleanFileExt, cleanFileExtFlagName, false, "Attempt to clean original file extension")
	sortCmd.Flags().BoolVar(&stopOnError, stopOnErrorFlagName, false, "Exit on first error")
	sortCmd.Flags().StringArrayVar(&fileExts,
		fileExtsFlagName,
		[]string{},
		fmt.Sprintf("Allowlist of file extensions to match on (default: %v)", mediasorter.DefaultFileTypes))

	_ = sortCmd.MarkFlagRequired(sourceDirFlagName)
	rootCmd.AddCommand(sortCmd)
}
