package humane

import (
	"fmt"
	"time"
)

// TimeFormatter formats one time relative to another the way Finder-
// adjacent tools do: symmetric "X ago" / "X from now" phrasing, with no
// "about" prefix on the hour bucket -- Swift's RelativeDateTimeFormatter
// has no such prefix either, and Go's justincampbell/timeago (which
// does add one) is exactly what this package replaces.
//
// The zero value has CollapseMinute set to false (second-level
// granularity for anything under a minute); use NewTimeFormatter for
// the collapsed default described below. Go structs can't default a
// bool field to true from a bare literal, so the constructor is the
// recommended way to get it -- TimeFormatter{CollapseMinute: true}
// works identically if you'd rather not use it.
type TimeFormatter struct {
	// CollapseMinute renders any duration under 60 seconds as "less
	// than a minute ago" / "less than a minute from now" instead of
	// counting seconds. Rails' distance_of_time_in_words, Go's
	// justincampbell/timeago, and zouk's own RelativeDateTimeFormatter
	// wrapper all do this in practice -- Swift's formatter has no such
	// bucket natively, so there's no "pure" behavior being overridden
	// here, just a convenience every real reference already reaches for.
	CollapseMinute bool
}

// NewTimeFormatter returns a TimeFormatter with CollapseMinute enabled,
// the recommended default.
func NewTimeFormatter() TimeFormatter {
	return TimeFormatter{CollapseMinute: true}
}

// Format returns t relative to relativeTo as a human-readable string.
//
//	f := NewTimeFormatter()
//	f.Format(t, t)                          == "less than a minute ago"
//	f.Format(t, t.Add(3*time.Minute))        == "3 minutes ago"
//	f.Format(t, t.Add(-3*time.Minute))       == "3 minutes from now"
//	f.Format(t, t.Add(15*time.Hour))         == "15 hours ago"
//	f.Format(t, t.Add(30*time.Hour))         == "1 day ago"
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
