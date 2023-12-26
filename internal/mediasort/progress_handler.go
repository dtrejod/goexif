package mediasort

import (
	"context"

	"github.com/dtrejod/goexif/internal/ilog"
	"go.uber.org/zap"
)

const (
	// logPercentage is the percentage increments at which we log progress
	logPercentage = 5 // log every 5%
)

type progressTracker struct {
	currentMediaIndex int
	logThreshold      int
	logNextThreshold  int
	totalMediaFiles   int
}

// handle is used to track progress of a scan run. It requires traversing
// through a folder once to accumulate a count the total number of files that
// will be handled. Then on second pass it will log occasionally the overall
// progress.
func (s *progressTracker) handle(ctx context.Context, isAccumulating bool) {
	if isAccumulating {
		s.totalMediaFiles++
		return
	}

	s.currentMediaIndex++
	if s.logThreshold == 0 {
		// calculate log threshold on 1st run
		s.logThreshold = int(float64(s.totalMediaFiles) * float64(logPercentage) / 100)
		s.logNextThreshold = s.logThreshold
	}

	// Log Progress every logPercent and at the very beginning and end
	if s.currentMediaIndex == 1 || s.currentMediaIndex == s.logNextThreshold || s.currentMediaIndex == s.totalMediaFiles {
		// double threshold every time we hit the threshold
		if s.currentMediaIndex == s.logNextThreshold {
			s.logNextThreshold += s.logThreshold
		}

		percentageComplete := int(float64(s.currentMediaIndex) / float64(s.totalMediaFiles) * 100)
		ilog.FromContext(ctx).Info("Current progress.",
			zap.Int("current", s.currentMediaIndex),
			zap.Int("total", s.totalMediaFiles),
			zap.Int("percentage", percentageComplete))
	}
}
