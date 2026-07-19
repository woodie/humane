package humane_test

import (
	"testing"
	"time"

	"github.com/sclevine/spec"
	. "github.com/woodie/expect"

	"github.com/woodie/humane"
)

// ptr is a small test-only helper -- DistanceInTime/TimeAgo take *time.Time
// so a nil at is expressible, but Go can't take the address of a literal or
// a func result inline.
func ptr(t time.Time) *time.Time { return &t }

func TestTime(t *testing.T) {
	spec.RunAliased(t, "Time", func(t *testing.T, describe, context spec.Describe, it spec.S, _, _ func(func())) {
		describe("DistanceInTime", func() {
			base := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)

			context("with no options (the recommended defaults: Approximate true, IncludeSeconds false -- matching ActionView's own defaults)", func() {
				context("just now", func() {
					it("displays less than a minute ago", func() {
						Expect(t, humane.DistanceInTime(ptr(base), base)).To(Equal("less than a minute ago"))
					})
				})

				context("45 seconds ago", func() {
					it("rounds up to 1 minute ago (past the 30-second cutoff)", func() {
						Expect(t, humane.DistanceInTime(ptr(base.Add(-45*time.Second)), base)).To(Equal("1 minute ago"))
					})
				})

				context("1 minute ago", func() {
					it("displays 1 minute ago, singular", func() {
						Expect(t, humane.DistanceInTime(ptr(base.Add(-1*time.Minute)), base)).To(Equal("1 minute ago"))
					})
				})

				context("3 minutes ago", func() {
					it("displays 3 minutes ago", func() {
						Expect(t, humane.DistanceInTime(ptr(base.Add(-3*time.Minute)), base)).To(Equal("3 minutes ago"))
					})
				})

				context("1 hour ago", func() {
					it("displays about 1 hour ago", func() {
						Expect(t, humane.DistanceInTime(ptr(base.Add(-1*time.Hour)), base)).To(Equal("about 1 hour ago"))
					})
				})

				context("15 hours ago", func() {
					it("displays about 15 hours ago", func() {
						Expect(t, humane.DistanceInTime(ptr(base.Add(-15*time.Hour)), base)).To(Equal("about 15 hours ago"))
					})
				})

				context("30 hours ago", func() {
					it("rolls up to 1 day ago, with no about (ActionView's table has none on the day bucket)", func() {
						Expect(t, humane.DistanceInTime(ptr(base.Add(-30*time.Hour)), base)).To(Equal("1 day ago"))
					})
				})

				context("3 days ago", func() {
					it("displays 3 days ago", func() {
						Expect(t, humane.DistanceInTime(ptr(base.Add(-3*24*time.Hour)), base)).To(Equal("3 days ago"))
					})
				})

				context("45 seconds from now", func() {
					it("rounds up to in 1 minute (past the 30-second cutoff)", func() {
						Expect(t, humane.DistanceInTime(ptr(base.Add(45*time.Second)), base)).To(Equal("in 1 minute"))
					})
				})

				context("3 minutes from now", func() {
					it("displays in 3 minutes", func() {
						Expect(t, humane.DistanceInTime(ptr(base.Add(3*time.Minute)), base)).To(Equal("in 3 minutes"))
					})
				})

				context("3 hours from now", func() {
					it("displays in about 3 hours", func() {
						Expect(t, humane.DistanceInTime(ptr(base.Add(3*time.Hour)), base)).To(Equal("in about 3 hours"))
					})
				})
			})

			context("with IncludeSeconds: true", func() {
				opts := humane.TimeOptions{IncludeSeconds: true}

				context("just now", func() {
					it("displays 0 seconds ago", func() {
						Expect(t, humane.DistanceInTime(ptr(base), base, opts)).To(Equal("0 seconds ago"))
					})
				})

				context("1 second ago", func() {
					it("displays 1 second ago, singular", func() {
						Expect(t, humane.DistanceInTime(ptr(base.Add(-1*time.Second)), base, opts)).To(Equal("1 second ago"))
					})
				})

				context("45 seconds ago", func() {
					it("displays 45 seconds ago", func() {
						Expect(t, humane.DistanceInTime(ptr(base.Add(-45*time.Second)), base, opts)).To(Equal("45 seconds ago"))
					})
				})

				context("45 seconds from now", func() {
					it("displays in 45 seconds", func() {
						Expect(t, humane.DistanceInTime(ptr(base.Add(45*time.Second)), base, opts)).To(Equal("in 45 seconds"))
					})
				})
			})

			context("with Approximate: false", func() {
				opts := humane.TimeOptions{Approximate: humane.Bool(false)}

				context("1 hour ago", func() {
					it("displays the exact count, no about prefix", func() {
						Expect(t, humane.DistanceInTime(ptr(base.Add(-1*time.Hour)), base, opts)).To(Equal("1 hour ago"))
					})
				})

				context("15 hours ago", func() {
					it("displays 15 hours ago", func() {
						Expect(t, humane.DistanceInTime(ptr(base.Add(-15*time.Hour)), base, opts)).To(Equal("15 hours ago"))
					})
				})
			})

			describe("nil handling", func() {
				context("when at is nil and WhenNil is set", func() {
					it("returns WhenNil without formatting", func() {
						opts := humane.TimeOptions{WhenNil: "an unknown time"}
						Expect(t, humane.DistanceInTime(nil, base, opts)).To(Equal("an unknown time"))
					})
				})

				context("when at is nil and WhenNil is left unset", func() {
					it("returns an empty string", func() {
						Expect(t, humane.DistanceInTime(nil, base)).To(Equal(""))
					})
				})
			})

			// Boundary regression coverage for the ActionView distance_of_time_in_words bucket
			// table this approximate-distance behavior ports, truncated at the "1 day" row
			// since month/year buckets are out of scope. Each pair straddles a cutoff second
			// from that table to lock in exactly where the wording flips.
			describe("at the approximate-distance bucket table boundaries", func() {
				context("with Approximate: false", func() {
					opts := humane.TimeOptions{Approximate: humane.Bool(false)}

					it("29s stays less than a minute, 30s rounds up to 1 minute", func() {
						Expect(t, humane.DistanceInTime(ptr(base.Add(-29*time.Second)), base, opts)).To(Equal("less than a minute ago"))
						Expect(t, humane.DistanceInTime(ptr(base.Add(-30*time.Second)), base, opts)).To(Equal("1 minute ago"))
					})

					it("89s stays 1 minute, 90s rounds up to 2 minutes", func() {
						Expect(t, humane.DistanceInTime(ptr(base.Add(-89*time.Second)), base, opts)).To(Equal("1 minute ago"))
						Expect(t, humane.DistanceInTime(ptr(base.Add(-90*time.Second)), base, opts)).To(Equal("2 minutes ago"))
					})

					it("44:29 stays 44 minutes, 44:30 rounds up to 1 hour", func() {
						Expect(t, humane.DistanceInTime(ptr(base.Add(-(44*time.Minute+29*time.Second))), base, opts)).To(Equal("44 minutes ago"))
						Expect(t, humane.DistanceInTime(ptr(base.Add(-(44*time.Minute+30*time.Second))), base, opts)).To(Equal("1 hour ago"))
					})

					it("89:29 stays 1 hour, 89:30 rounds up to 2 hours", func() {
						Expect(t, humane.DistanceInTime(ptr(base.Add(-(89*time.Minute+29*time.Second))), base, opts)).To(Equal("1 hour ago"))
						Expect(t, humane.DistanceInTime(ptr(base.Add(-(89*time.Minute+30*time.Second))), base, opts)).To(Equal("2 hours ago"))
					})

					it("23:59:29 stays 24 hours, 23:59:30 rounds up to 1 day", func() {
						Expect(t, humane.DistanceInTime(ptr(base.Add(-(23*time.Hour+59*time.Minute+29*time.Second))), base, opts)).To(Equal("24 hours ago"))
						Expect(t, humane.DistanceInTime(ptr(base.Add(-(23*time.Hour+59*time.Minute+30*time.Second))), base, opts)).To(Equal("1 day ago"))
					})
				})

				context("with no options (Approximate true by default)", func() {
					it("44:29 has no about, 44:30 gains about (entering the hour bucket)", func() {
						Expect(t, humane.DistanceInTime(ptr(base.Add(-(44*time.Minute+29*time.Second))), base)).To(Equal("44 minutes ago"))
						Expect(t, humane.DistanceInTime(ptr(base.Add(-(44*time.Minute+30*time.Second))), base)).To(Equal("about 1 hour ago"))
					})

					it("23:59:29 keeps about, 23:59:30 drops about (entering the day bucket)", func() {
						Expect(t, humane.DistanceInTime(ptr(base.Add(-(23*time.Hour+59*time.Minute+29*time.Second))), base)).To(Equal("about 24 hours ago"))
						Expect(t, humane.DistanceInTime(ptr(base.Add(-(23*time.Hour+59*time.Minute+30*time.Second))), base)).To(Equal("1 day ago"))
					})
				})
			})
		})

		// TimeAgo is a thin one-argument convenience over DistanceInTime, supplying
		// time.Now() as relativeTo -- see DistanceInTime above for the exhaustive
		// wording/bucket coverage this doesn't need to repeat.
		describe("TimeAgo", func() {
			context("just now", func() {
				it("displays less than a minute ago", func() {
					Expect(t, humane.TimeAgo(time.Now())).To(Equal("less than a minute ago"))
				})
			})

			context("3 minutes ago", func() {
				it("forwards to DistanceInTime with time.Now() as relativeTo", func() {
					Expect(t, humane.TimeAgo(time.Now().Add(-3*time.Minute))).To(Equal("3 minutes ago"))
				})
			})
		})
	})
}
