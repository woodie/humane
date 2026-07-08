package humane

import (
	"testing"
	"time"
)

func TestTimeFormatter_Format_collapseMinute(t *testing.T) {
	base := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)
	f := NewTimeFormatter() // CollapseMinute: true

	cases := []struct {
		name string
		when time.Time
		want string
	}{
		{"just now", base, "less than a minute ago"},
		{"45 seconds ago", base.Add(-45 * time.Second), "less than a minute ago"},
		{"3 minutes ago", base.Add(-3 * time.Minute), "3 minutes ago"},
		{"1 minute ago (singular)", base.Add(-1 * time.Minute), "1 minute ago"},
		{"15 hours ago, no 'about' prefix", base.Add(-15 * time.Hour), "15 hours ago"},
		{"1 hour ago (singular)", base.Add(-1 * time.Hour), "1 hour ago"},
		{"30 hours ago rolls to 1 day", base.Add(-30 * time.Hour), "1 day ago"},
		{"3 days ago", base.Add(-72 * time.Hour), "3 days ago"},
		{"45 seconds from now", base.Add(45 * time.Second), "less than a minute from now"},
		{"3 minutes from now", base.Add(3 * time.Minute), "3 minutes from now"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := f.Format(c.when, base); got != c.want {
				t.Errorf("Format(%v, base) = %q, want %q", c.when, got, c.want)
			}
		})
	}
}

func TestTimeFormatter_Format_precise(t *testing.T) {
	base := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)
	f := TimeFormatter{} // CollapseMinute: false -- zero value, second-level detail

	cases := []struct {
		name string
		when time.Time
		want string
	}{
		{"just now", base, "0 seconds ago"},
		{"45 seconds ago", base.Add(-45 * time.Second), "45 seconds ago"},
		{"1 second ago (singular)", base.Add(-1 * time.Second), "1 second ago"},
		{"45 seconds from now", base.Add(45 * time.Second), "45 seconds from now"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := f.Format(c.when, base); got != c.want {
				t.Errorf("Format(%v, base) = %q, want %q", c.when, got, c.want)
			}
		})
	}
}
