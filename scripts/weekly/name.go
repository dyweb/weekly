package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/now"
)

// name.go generates weekly file name based on create time and recent issues.
// See https://github.com/dyweb/weekly/issues/186 for background.

// All the issues that are created within a week are aligned to that week's Monday.
// The comments are merged into a single file.

func FileName(createTime time.Time) string {
	// Copied from https://github.com/dyweb/dy-bot/blob/2fedb230d6ba21f0ebb9f1091f27d8482c331772/pkg/util/weeklyutil/weekly.go#L11
	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	beijing := time.FixedZone("Beijing Time", secondsEastOfUTC)

	cfg := now.Config{
		WeekStartDay: time.Monday,
		TimeLocation: beijing,
	}
	date := cfg.With(createTime).BeginningOfWeek()
	return fmt.Sprintf("%d/%d-%.2d-%.2d-weekly.md", date.Year(), date.Year(), int(date.Month()), date.Day())
}
