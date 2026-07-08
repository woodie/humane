package humane

import (
	"fmt"
	"time"
)

// TimeFormatter formats one time relative to another the way Finder-adjacent tools do.
type TimeFormatter struct {
	// CollapseMinute buckets anything under a minute as "less than a minute ago/from now".
	CollapseMinute bool
}

// NewTimeFormatter returns a TimeFormatter with CollapseMinute enabled, the recommended default.
func NewTimeFormatter() TimeFormatter {
	return TimeFormatter{CollapseMinute: true}
}

// Format returns t relative to relativeTo as a human-readable string.
func (f TimeFormatter) Format(t, relativeTo time.Time) string {
	d := relativeTo.Sub(t)
	future := d < 0
	if future {
		d = -d
	}

	var text string
	switch {
	case f.CollapseMinute && d < time.Minute:
		if future {
			return "less than a minute from now"
		}
		return "less than a minute ago"
	case d < time.Minute:
		text = pluralize(int(d.Seconds()), "second")
	case d < time.Hour:
		text = pluralize(int(d.Minutes()+0.5), "minute")
	case d < 24*time.Hour:
		text = pluralize(int(d.Hours()+0.5), "hour")
	default:
		text = pluralize(int(d.Hours()/24+0.5), "day")
	}

	if future {
		return text + " from now"
	}
	return text + " ago"
}

func pluralize(n int, unit string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", unit)
	}
	return fmt.Sprintf("%d %ss", n, unit)
}
