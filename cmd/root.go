package cmd

import (
	"context"

	"github.com/dtrejod/goexif/internal/ilog"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ctx         context.Context
	debug       bool
	logEncoding string
	rootCmd     = &cobra.Command{
		Use:               "goexif",
		Short:             "A tool for interacting with media files via their exif/file metadata",
		PersistentPreRunE: initLoggers,
	}
)

func initLoggers(_ *cobra.Command, _ []string) error {
	ctx = context.Background()
	var err error
	logConfig := zap.NewProductionConfig()
	if debug {
		logConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	logConfig.Encoding = logEncoding

	// if console encoding, make the logs more human readable
	if logEncoding == "console" {
		logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logConfig.EncoderConfig.ConsoleSeparator = " "
		logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	logger, err := logConfig.Build()
	if err != nil {
		return err
	}
	ctx = ilog.WithLogger(ctx, logger)
	return nil
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "run in debug mode")
	rootCmd.PersistentFlags().StringVarP(&logEncoding, "log-encoding", "e", "console", "log encoding (json or console)")
}

// Execute runs the root command tree
func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		ilog.FromContext(ctx).Error("goexif error", zap.Error(err))
		return 1
	}
	return 0
}
