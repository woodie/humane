package humane

import (
	"fmt"
	"time"
)

// TimeOptions configures TimeAgo. Omit entirely (or pass zero options) for
// the recommended defaults, which match ActionView's own: Approximate true,
// IncludeSeconds false. See docs/COMMENTS.md for why Approximate is a *bool
// rather than a bool.
type TimeOptions struct {
	IncludeSeconds bool
	Approximate    *bool
	WhenNil        string
}

// Bool returns a pointer to b, for use with TimeOptions.Approximate.
func Bool(b bool) *bool { return &b }

// TimeAgo formats at relative to relativeTo as a human-readable string,
// worded "X ago"/"in X" -- direction-aware, so the caller never has to know
// ahead of time whether at is in the past or future. If at is nil, returns
// opts.WhenNil without formatting; see docs/COMMENTS.md.
func TimeAgo(at *time.Time, relativeTo time.Time, opts ...TimeOptions) string {
	o := TimeOptions{}
	if len(opts) > 0 {
		o = opts[0]
	}

	if at == nil {
		return o.WhenNil
	}

	approximate := true
	if o.Approximate != nil {
		approximate = *o.Approximate
	}

	d := relativeTo.Sub(*at)
	future := d < 0
	if future {
		d = -d
	}

	if !o.IncludeSeconds && d < 30*time.Second {
		if future {
			return "in less than a minute"
		}
		return "less than a minute ago"
	}

	if o.IncludeSeconds && d < time.Minute {
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

	if approximate && approximable {
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
