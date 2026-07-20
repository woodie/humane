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
	spec.Run(t, "Time", func(t *testing.T, describe spec.G, it spec.S) {
		context, before := describe, it.Before

		describe("DistanceInTime", func() {
			var at *time.Time
			var base time.Time
			var opts humane.TimeOptions
			subject := func() string { return humane.DistanceInTime(at, base, opts) }

			before(func() {
				base = time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)
				opts = humane.TimeOptions{}
			})

			context("just now", func() {
				before(func() { at = ptr(base) })

				context("with no options (the recommended defaults: Approximate true, IncludeSeconds false -- matching ActionView's own defaults)", func() {
					it("displays less than a minute ago", func() {
						expect(subject(), t).To(Equal("less than a minute ago"))
					})
				})

				context("with IncludeSeconds: true", func() {
					before(func() { opts = humane.TimeOptions{IncludeSeconds: true} })

					it("displays 0 seconds ago", func() {
						expect(subject(), t).To(Equal("0 seconds ago"))
					})
				})
			})

			context("1 second ago", func() {
				before(func() { at = ptr(base.Add(-1 * time.Second)) })

				context("with IncludeSeconds: true", func() {
					before(func() { opts = humane.TimeOptions{IncludeSeconds: true} })

					it("displays 1 second ago, singular", func() {
						expect(subject(), t).To(Equal("1 second ago"))
					})
				})
			})

			context("45 seconds ago", func() {
				before(func() { at = ptr(base.Add(-45 * time.Second)) })

				context("with no options", func() {
					it("rounds up to 1 minute ago (past the 30-second cutoff)", func() {
						expect(subject(), t).To(Equal("1 minute ago"))
					})
				})

				context("with IncludeSeconds: true", func() {
					before(func() { opts = humane.TimeOptions{IncludeSeconds: true} })

					it("displays 45 seconds ago", func() {
						expect(subject(), t).To(Equal("45 seconds ago"))
					})
				})
			})

			context("1 minute ago", func() {
				before(func() { at = ptr(base.Add(-1 * time.Minute)) })

				context("with no options", func() {
					it("displays 1 minute ago, singular", func() {
						expect(subject(), t).To(Equal("1 minute ago"))
					})
				})
			})

			context("3 minutes ago", func() {
				before(func() { at = ptr(base.Add(-3 * time.Minute)) })

				context("with no options", func() {
					it("displays 3 minutes ago", func() {
						expect(subject(), t).To(Equal("3 minutes ago"))
					})
				})
			})

			context("1 hour ago", func() {
				before(func() { at = ptr(base.Add(-1 * time.Hour)) })

				context("with no options", func() {
					it("displays about 1 hour ago", func() {
						expect(subject(), t).To(Equal("about 1 hour ago"))
					})
				})

				context("with Approximate: false", func() {
					before(func() { opts = humane.TimeOptions{Approximate: humane.Bool(false)} })

					it("displays the exact count, no about prefix", func() {
						expect(subject(), t).To(Equal("1 hour ago"))
					})
				})
			})

			context("15 hours ago", func() {
				before(func() { at = ptr(base.Add(-15 * time.Hour)) })

				context("with no options", func() {
					it("displays about 15 hours ago", func() {
						expect(subject(), t).To(Equal("about 15 hours ago"))
					})
				})

				context("with Approximate: false", func() {
					before(func() { opts = humane.TimeOptions{Approximate: humane.Bool(false)} })

					it("displays 15 hours ago", func() {
						expect(subject(), t).To(Equal("15 hours ago"))
					})
				})
			})

			context("30 hours ago", func() {
				before(func() { at = ptr(base.Add(-30 * time.Hour)) })

				context("with no options", func() {
					it("rolls up to 1 day ago, with no about (ActionView's table has none on the day bucket)", func() {
						expect(subject(), t).To(Equal("1 day ago"))
					})
				})
			})

			context("3 days ago", func() {
				before(func() { at = ptr(base.Add(-3 * 24 * time.Hour)) })

				context("with no options", func() {
					it("displays 3 days ago", func() {
						expect(subject(), t).To(Equal("3 days ago"))
					})
				})
			})

			context("45 seconds from now", func() {
				before(func() { at = ptr(base.Add(45 * time.Second)) })

				context("with no options", func() {
					it("rounds up to in 1 minute (past the 30-second cutoff)", func() {
						expect(subject(), t).To(Equal("in 1 minute"))
					})
				})

				context("with IncludeSeconds: true", func() {
					before(func() { opts = humane.TimeOptions{IncludeSeconds: true} })

					it("displays in 45 seconds", func() {
						expect(subject(), t).To(Equal("in 45 seconds"))
					})
				})
			})

			context("3 minutes from now", func() {
				before(func() { at = ptr(base.Add(3 * time.Minute)) })

				context("with no options", func() {
					it("displays in 3 minutes", func() {
						expect(subject(), t).To(Equal("in 3 minutes"))
					})
				})
			})

			context("3 hours from now", func() {
				before(func() { at = ptr(base.Add(3 * time.Hour)) })

				context("with no options", func() {
					it("displays in about 3 hours", func() {
						expect(subject(), t).To(Equal("in about 3 hours"))
					})
				})
			})

			describe("nil handling", func() {
				before(func() { at = nil })

				context("when at is nil and WhenNil is set", func() {
					before(func() { opts = humane.TimeOptions{WhenNil: "an unknown time"} })

					it("returns WhenNil without formatting", func() {
						expect(subject(), t).To(Equal("an unknown time"))
					})
				})

				context("when at is nil and WhenNil is left unset", func() {
					it("returns an empty string", func() {
						expect(subject(), t).To(Equal(""))
					})
				})
			})

			// Boundary regression coverage for ActionView's distance_of_time_in_words bucket table (truncated at the "1 day" row); each context below sits on one cutoff second from that table.
			describe("at the approximate-distance bucket table boundaries", func() {
				context("with Approximate: false", func() {
					before(func() { opts = humane.TimeOptions{Approximate: humane.Bool(false)} })

					context("29 seconds ago", func() {
						before(func() { at = ptr(base.Add(-29 * time.Second)) })

						it("stays less than a minute", func() {
							expect(subject(), t).To(Equal("less than a minute ago"))
						})
					})

					context("30 seconds ago", func() {
						before(func() { at = ptr(base.Add(-30 * time.Second)) })

						it("rounds up to 1 minute", func() {
							expect(subject(), t).To(Equal("1 minute ago"))
						})
					})

					context("89 seconds ago", func() {
						before(func() { at = ptr(base.Add(-89 * time.Second)) })

						it("stays 1 minute", func() {
							expect(subject(), t).To(Equal("1 minute ago"))
						})
					})

					context("90 seconds ago", func() {
						before(func() { at = ptr(base.Add(-90 * time.Second)) })

						it("rounds up to 2 minutes", func() {
							expect(subject(), t).To(Equal("2 minutes ago"))
						})
					})

					context("44 minutes 29 seconds ago", func() {
						before(func() { at = ptr(base.Add(-(44*time.Minute + 29*time.Second))) })

						it("stays 44 minutes", func() {
							expect(subject(), t).To(Equal("44 minutes ago"))
						})
					})

					context("44 minutes 30 seconds ago", func() {
						before(func() { at = ptr(base.Add(-(44*time.Minute + 30*time.Second))) })

						it("rounds up to 1 hour", func() {
							expect(subject(), t).To(Equal("1 hour ago"))
						})
					})

					context("89 minutes 29 seconds ago", func() {
						before(func() { at = ptr(base.Add(-(89*time.Minute + 29*time.Second))) })

						it("stays 1 hour", func() {
							expect(subject(), t).To(Equal("1 hour ago"))
						})
					})

					context("89 minutes 30 seconds ago", func() {
						before(func() { at = ptr(base.Add(-(89*time.Minute + 30*time.Second))) })

						it("rounds up to 2 hours", func() {
							expect(subject(), t).To(Equal("2 hours ago"))
						})
					})

					context("23 hours 59 minutes 29 seconds ago", func() {
						before(func() { at = ptr(base.Add(-(23*time.Hour + 59*time.Minute + 29*time.Second))) })

						it("stays 24 hours", func() {
							expect(subject(), t).To(Equal("24 hours ago"))
						})
					})

					context("23 hours 59 minutes 30 seconds ago", func() {
						before(func() { at = ptr(base.Add(-(23*time.Hour + 59*time.Minute + 30*time.Second))) })

						it("rounds up to 1 day", func() {
							expect(subject(), t).To(Equal("1 day ago"))
						})
					})
				})

				context("with no options (Approximate true by default)", func() {
					context("44 minutes 29 seconds ago", func() {
						before(func() { at = ptr(base.Add(-(44*time.Minute + 29*time.Second))) })

						it("has no about", func() {
							expect(subject(), t).To(Equal("44 minutes ago"))
						})
					})

					context("44 minutes 30 seconds ago", func() {
						before(func() { at = ptr(base.Add(-(44*time.Minute + 30*time.Second))) })

						it("gains about, entering the hour bucket", func() {
							expect(subject(), t).To(Equal("about 1 hour ago"))
						})
					})

					context("23 hours 59 minutes 29 seconds ago", func() {
						before(func() { at = ptr(base.Add(-(23*time.Hour + 59*time.Minute + 29*time.Second))) })

						it("keeps about", func() {
							expect(subject(), t).To(Equal("about 24 hours ago"))
						})
					})

					context("23 hours 59 minutes 30 seconds ago", func() {
						before(func() { at = ptr(base.Add(-(23*time.Hour + 59*time.Minute + 30*time.Second))) })

						it("drops about, entering the day bucket", func() {
							expect(subject(), t).To(Equal("1 day ago"))
						})
					})
				})
			})
		})

		// TimeAgo is a thin one-argument convenience over DistanceInTime, supplying
		// time.Now() as relativeTo -- see DistanceInTime above for the exhaustive
		// wording/bucket coverage this doesn't need to repeat.
		describe("TimeAgo", func() {
			var when time.Time
			subject := func() string { return humane.TimeAgo(when) }

			context("just now", func() {
				before(func() { when = time.Now() })

				it("displays less than a minute ago", func() {
					expect(subject(), t).To(Equal("less than a minute ago"))
				})
			})

			context("3 minutes ago", func() {
				before(func() { when = time.Now().Add(-3 * time.Minute) })

				it("forwards to DistanceInTime with time.Now() as relativeTo", func() {
					expect(subject(), t).To(Equal("3 minutes ago"))
				})
			})
		})
	})
}
