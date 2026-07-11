package humane

import (
	"fmt"
	"time"
)

// TimeFormatter formats one time relative to another the way Finder-adjacent tools do.
type TimeFormatter struct {
	// IncludeSeconds shows exact seconds under a minute instead of collapsing to "less than a minute ago/in less than a minute". Zero value (false) matches ActionView's include_seconds default.
	IncludeSeconds bool
	// Approximate prefixes "about"/"in about" on the hour-scale buckets (1 hour, and 2..24 hours), matching ActionView's distance_of_time_in_words wording for those buckets. Zero value (false) matches Foundation's raw output. See docs/COMMENTS.md and humane-ruby issue #1 for the full bucket table this ports.
	Approximate bool
}

// NewTimeFormatter returns a TimeFormatter with the recommended default -- now identical to the zero value, since IncludeSeconds' zero value (false) already is that default.
func NewTimeFormatter() TimeFormatter {
	return TimeFormatter{}
}

// Format returns at relative to relativeTo as a human-readable string.
func (f TimeFormatter) Format(at, relativeTo time.Time) string {
	d := relativeTo.Sub(at)
	future := d < 0
	if future {
		d = -d
	}

	if !f.IncludeSeconds && d < 30*time.Second {
		if future {
			return "in less than a minute"
		}
		return "less than a minute ago"
	}

	if f.IncludeSeconds && d < time.Minute {
		return wrap(pluralize(int(d.Seconds()), "second"), future)
	}

	// Buckets come from distanceInMinutes, not raw seconds re-divided per unit -- see docs/COMMENTS.md.
	distanceInMinutes := int(d.Minutes() + 0.5)

	var text string
	var approximable bool
	switch {
	case distanceInMinutes == 1:
		text = "1 minute"
	case distanceInMinutes <= 44:
		text = pluralize(distanceInMinutes, "minute")
	case distanceInMinutes <= 89:
		text = "1 hour"
		approximable = true
	case distanceInMinutes <= 1439:
		text = pluralize(int(float64(distanceInMinutes)/60.0+0.5), "hour")
		approximable = true
	case distanceInMinutes <= 2519:
		text = "1 day"
	default:
		text = pluralize(int(float64(distanceInMinutes)/1440.0+0.5), "day")
	}

	if f.Approximate && approximable {
		text = "about " + text
	}

	return wrap(text, future)
}

func wrap(text string, future bool) string {
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
