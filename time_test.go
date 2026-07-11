package humane_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/woodie/humane"
)

// ptr is a small test-only helper -- TimeAgo takes *time.Time so a nil at is
// expressible, but Go can't take the address of a literal or a func result
// inline.
func ptr(t time.Time) *time.Time { return &t }

var _ = Describe("TimeAgo", func() {
	base := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)

	Describe("with no options (the recommended defaults: Approximate true, IncludeSeconds false -- matching ActionView's own defaults)", func() {
		Context("just now", func() {
			It("displays less than a minute ago", func() {
				Expect(humane.TimeAgo(ptr(base), base)).To(Equal("less than a minute ago"))
			})
		})

		Context("45 seconds ago", func() {
			It("rounds up to 1 minute ago (past the 30-second cutoff)", func() {
				Expect(humane.TimeAgo(ptr(base.Add(-45*time.Second)), base)).To(Equal("1 minute ago"))
			})
		})

		Context("1 minute ago", func() {
			It("displays 1 minute ago, singular", func() {
				Expect(humane.TimeAgo(ptr(base.Add(-1*time.Minute)), base)).To(Equal("1 minute ago"))
			})
		})

		Context("3 minutes ago", func() {
			It("displays 3 minutes ago", func() {
				Expect(humane.TimeAgo(ptr(base.Add(-3*time.Minute)), base)).To(Equal("3 minutes ago"))
			})
		})

		Context("1 hour ago", func() {
			It("displays about 1 hour ago", func() {
				Expect(humane.TimeAgo(ptr(base.Add(-1*time.Hour)), base)).To(Equal("about 1 hour ago"))
			})
		})

		Context("15 hours ago", func() {
			It("displays about 15 hours ago", func() {
				Expect(humane.TimeAgo(ptr(base.Add(-15*time.Hour)), base)).To(Equal("about 15 hours ago"))
			})
		})

		Context("30 hours ago", func() {
			It("rolls up to 1 day ago, with no about (ActionView's table has none on the day bucket)", func() {
				Expect(humane.TimeAgo(ptr(base.Add(-30*time.Hour)), base)).To(Equal("1 day ago"))
			})
		})

		Context("3 days ago", func() {
			It("displays 3 days ago", func() {
				Expect(humane.TimeAgo(ptr(base.Add(-3*24*time.Hour)), base)).To(Equal("3 days ago"))
			})
		})

		Context("45 seconds from now", func() {
			It("rounds up to in 1 minute (past the 30-second cutoff)", func() {
				Expect(humane.TimeAgo(ptr(base.Add(45*time.Second)), base)).To(Equal("in 1 minute"))
			})
		})

		Context("3 minutes from now", func() {
			It("displays in 3 minutes", func() {
				Expect(humane.TimeAgo(ptr(base.Add(3*time.Minute)), base)).To(Equal("in 3 minutes"))
			})
		})

		Context("3 hours from now", func() {
			It("displays in about 3 hours", func() {
				Expect(humane.TimeAgo(ptr(base.Add(3*time.Hour)), base)).To(Equal("in about 3 hours"))
			})
		})
	})

	Describe("with IncludeSeconds: true", func() {
		opts := humane.TimeOptions{IncludeSeconds: true}

		Context("just now", func() {
			It("displays 0 seconds ago", func() {
				Expect(humane.TimeAgo(ptr(base), base, opts)).To(Equal("0 seconds ago"))
			})
		})

		Context("1 second ago", func() {
			It("displays 1 second ago, singular", func() {
				Expect(humane.TimeAgo(ptr(base.Add(-1*time.Second)), base, opts)).To(Equal("1 second ago"))
			})
		})

		Context("45 seconds ago", func() {
			It("displays 45 seconds ago", func() {
				Expect(humane.TimeAgo(ptr(base.Add(-45*time.Second)), base, opts)).To(Equal("45 seconds ago"))
			})
		})

		Context("45 seconds from now", func() {
			It("displays in 45 seconds", func() {
				Expect(humane.TimeAgo(ptr(base.Add(45*time.Second)), base, opts)).To(Equal("in 45 seconds"))
			})
		})
	})

	Describe("with Approximate: false", func() {
		opts := humane.TimeOptions{Approximate: humane.Bool(false)}

		Context("1 hour ago", func() {
			It("displays the exact count, no about prefix", func() {
				Expect(humane.TimeAgo(ptr(base.Add(-1*time.Hour)), base, opts)).To(Equal("1 hour ago"))
			})
		})

		Context("15 hours ago", func() {
			It("displays 15 hours ago", func() {
				Expect(humane.TimeAgo(ptr(base.Add(-15*time.Hour)), base, opts)).To(Equal("15 hours ago"))
			})
		})
	})

	Describe("nil handling", func() {
		Context("when at is nil and WhenNil is set", func() {
			It("returns WhenNil without formatting", func() {
				opts := humane.TimeOptions{WhenNil: "an unknown time"}
				Expect(humane.TimeAgo(nil, base, opts)).To(Equal("an unknown time"))
			})
		})

		Context("when at is nil and WhenNil is left unset", func() {
			It("returns an empty string", func() {
				Expect(humane.TimeAgo(nil, base)).To(Equal(""))
			})
		})
	})

	// Boundary regression coverage for the ActionView distance_of_time_in_words bucket
	// table this approximate-distance behavior ports, truncated at the "1 day" row
	// since month/year buckets are out of scope. Each pair straddles a cutoff second
	// from that table to lock in exactly where the wording flips.
	Describe("at the approximate-distance bucket table boundaries", func() {
		Context("with Approximate: false", func() {
			opts := humane.TimeOptions{Approximate: humane.Bool(false)}

			It("29s stays less than a minute, 30s rounds up to 1 minute", func() {
				Expect(humane.TimeAgo(ptr(base.Add(-29*time.Second)), base, opts)).To(Equal("less than a minute ago"))
				Expect(humane.TimeAgo(ptr(base.Add(-30*time.Second)), base, opts)).To(Equal("1 minute ago"))
			})

			It("89s stays 1 minute, 90s rounds up to 2 minutes", func() {
				Expect(humane.TimeAgo(ptr(base.Add(-89*time.Second)), base, opts)).To(Equal("1 minute ago"))
				Expect(humane.TimeAgo(ptr(base.Add(-90*time.Second)), base, opts)).To(Equal("2 minutes ago"))
			})

			It("44:29 stays 44 minutes, 44:30 rounds up to 1 hour", func() {
				Expect(humane.TimeAgo(ptr(base.Add(-(44*time.Minute+29*time.Second))), base, opts)).To(Equal("44 minutes ago"))
				Expect(humane.TimeAgo(ptr(base.Add(-(44*time.Minute+30*time.Second))), base, opts)).To(Equal("1 hour ago"))
			})

			It("89:29 stays 1 hour, 89:30 rounds up to 2 hours", func() {
				Expect(humane.TimeAgo(ptr(base.Add(-(89*time.Minute+29*time.Second))), base, opts)).To(Equal("1 hour ago"))
				Expect(humane.TimeAgo(ptr(base.Add(-(89*time.Minute+30*time.Second))), base, opts)).To(Equal("2 hours ago"))
			})

			It("23:59:29 stays 24 hours, 23:59:30 rounds up to 1 day", func() {
				Expect(humane.TimeAgo(ptr(base.Add(-(23*time.Hour+59*time.Minute+29*time.Second))), base, opts)).To(Equal("24 hours ago"))
				Expect(humane.TimeAgo(ptr(base.Add(-(23*time.Hour+59*time.Minute+30*time.Second))), base, opts)).To(Equal("1 day ago"))
			})
		})

		Context("with no options (Approximate true by default)", func() {
			It("44:29 has no about, 44:30 gains about (entering the hour bucket)", func() {
				Expect(humane.TimeAgo(ptr(base.Add(-(44*time.Minute+29*time.Second))), base)).To(Equal("44 minutes ago"))
				Expect(humane.TimeAgo(ptr(base.Add(-(44*time.Minute+30*time.Second))), base)).To(Equal("about 1 hour ago"))
			})

			It("23:59:29 keeps about, 23:59:30 drops about (entering the day bucket)", func() {
				Expect(humane.TimeAgo(ptr(base.Add(-(23*time.Hour+59*time.Minute+29*time.Second))), base)).To(Equal("about 24 hours ago"))
				Expect(humane.TimeAgo(ptr(base.Add(-(23*time.Hour+59*time.Minute+30*time.Second))), base)).To(Equal("1 day ago"))
			})
		})
	})
})
