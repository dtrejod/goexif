package cmd

import (
	"os"

	"github.com/dtrejod/goexif/internal/ilog"
	"github.com/dtrejod/goexif/internal/mediadate"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	sourceFileFlagName = "src-file"
)

var (
	sourceFile string
)

var dateCmd = &cobra.Command{
	Use:   "date",
	Short: "Date prints the metadata datetime from the provided media file",
	Run:   dateRun,
}

func dateRun(_ *cobra.Command, _ []string) {
	if err := mediadate.Print(ctx, sourceFile, magicSignatureIn); err != nil {
		ilog.FromContext(ctx).Error("Failed to print date for mediafile",
			zap.String("sourceFile", sourceFile),
			zap.Error(err))
		os.Exit(1)
	}
}

func init() {
	dateCmd.Flags().StringVar(&sourceFile,
		sourceFileFlagName,
		"",
		"Source file to scan")
	dateCmd.Flags().BoolVar(&magicSignatureIn,
		magicSignatureInFlagName,
		false,
		"Ignore existing file extension and use magic signature instead when identifying files")

	_ = dateCmd.MarkFlagRequired(sourceFileFlagName)
	rootCmd.AddCommand(dateCmd)
}
