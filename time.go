package humane

import (
	"fmt"
	"time"
)

// TimeFormatter formats one time relative to another the way Finder-adjacent tools do.
type TimeFormatter struct {
	// IncludeSeconds shows exact seconds under a minute instead of collapsing to "less than a minute ago/in less than a minute". Zero value (false) matches ActionView's include_seconds default.
	IncludeSeconds bool
}

// NewTimeFormatter returns a TimeFormatter with the recommended default -- now identical to the zero value, since IncludeSeconds' zero value (false) already is that default.
func NewTimeFormatter() TimeFormatter {
	return TimeFormatter{}
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
	case !f.IncludeSeconds && d < time.Minute:
		if future {
			return "in less than a minute"
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
		return "in " + text
	}
	return text + " ago"
}

func pluralize(n int, unit string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", unit)
	}
	return fmt.Sprintf("%d %ss", n, unit)
}
