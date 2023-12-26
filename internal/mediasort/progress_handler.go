package mediasort

import (
	"context"

	"github.com/dtrejod/goexif/internal/ilog"
	"go.uber.org/zap"
)

type progressTracker struct {
	currentMediaIndex int
	totalMediaFiles   int
}

// handle is used to track progress of a scan run. It requires traversing through a folder once to accumulate a count the total number of files that will be handled. Then on second pass it will log occasionally the overall progress.
func (s *progressTracker) handle(ctx context.Context, isAccumulating bool) {
	if isAccumulating {
		s.totalMediaFiles++
		return
	}

	s.currentMediaIndex++

	percentageComplete := int(float64(s.currentMediaIndex) / float64(s.totalMediaFiles) * 100.0)
	ilog.FromContext(ctx).Debug("Current progress.",
		zap.Int("current", s.currentMediaIndex),
		zap.Int("total", s.totalMediaFiles),
		zap.Int("percentage", percentageComplete))

	// Log Progress
	if s.currentMediaIndex == 1 || (percentageComplete != 0 && percentageComplete%5 == 0) {
		ilog.FromContext(ctx).Info("Current progress.",
			zap.Int("current", s.currentMediaIndex),
			zap.Int("total", s.totalMediaFiles),
			zap.Int("percentage", percentageComplete))
	}
	return
}
