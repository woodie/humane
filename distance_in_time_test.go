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

func TestDistanceInTime(t *testing.T) {
	spec.Run(t, "humane.DistanceInTime", func(t *testing.T, context spec.G, it spec.S) {
		describe := context

		base := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)

		describe("with no options", func() {
			var at *time.Time
			var result string
			it.JustBeforeEach(func() { result = humane.DistanceInTime(at, base, humane.TimeOptions{}) })

			context("just now", func() {
				it.BeforeEach(func() { at = ptr(base) })

				it("displays less than a minute ago", func() {
					expect(result, t).To(Equal("less than a minute ago"))
				})
			})

			context("45 seconds ago", func() {
				it.BeforeEach(func() { at = ptr(base.Add(-45 * time.Second)) })

				it("rounds up to 1 minute ago (past the 30-second cutoff)", func() {
					expect(result, t).To(Equal("1 minute ago"))
				})
			})

			context("1 minute ago", func() {
				it.BeforeEach(func() { at = ptr(base.Add(-1 * time.Minute)) })

				it("displays 1 minute ago, singular", func() {
					expect(result, t).To(Equal("1 minute ago"))
				})
			})

			context("3 minutes ago", func() {
				it.BeforeEach(func() { at = ptr(base.Add(-3 * time.Minute)) })

				it("displays 3 minutes ago", func() {
					expect(result, t).To(Equal("3 minutes ago"))
				})
			})

			context("1 hour ago", func() {
				it.BeforeEach(func() { at = ptr(base.Add(-1 * time.Hour)) })

				it("displays about 1 hour ago", func() {
					expect(result, t).To(Equal("about 1 hour ago"))
				})
			})

			context("15 hours ago", func() {
				it.BeforeEach(func() { at = ptr(base.Add(-15 * time.Hour)) })

				it("displays about 15 hours ago", func() {
					expect(result, t).To(Equal("about 15 hours ago"))
				})
			})

			context("30 hours ago", func() {
				it.BeforeEach(func() { at = ptr(base.Add(-30 * time.Hour)) })

				it("rolls up to 1 day ago, with no about", func() {
					expect(result, t).To(Equal("1 day ago"))
				})
			})

			context("3 days ago", func() {
				it.BeforeEach(func() { at = ptr(base.Add(-3 * 24 * time.Hour)) })

				it("displays 3 days ago", func() {
					expect(result, t).To(Equal("3 days ago"))
				})
			})

			context("45 seconds from now", func() {
				it.BeforeEach(func() { at = ptr(base.Add(45 * time.Second)) })

				it("rounds up to in 1 minute (past the 30-second cutoff)", func() {
					expect(result, t).To(Equal("in 1 minute"))
				})
			})

			context("3 minutes from now", func() {
				it.BeforeEach(func() { at = ptr(base.Add(3 * time.Minute)) })

				it("displays in 3 minutes", func() {
					expect(result, t).To(Equal("in 3 minutes"))
				})
			})

			context("3 hours from now", func() {
				it.BeforeEach(func() { at = ptr(base.Add(3 * time.Hour)) })

				it("displays in about 3 hours", func() {
					expect(result, t).To(Equal("in about 3 hours"))
				})
			})
		})

		describe("with IncludeSeconds: true", func() {
			var at *time.Time
			var result string
			it.JustBeforeEach(func() {
				result = humane.DistanceInTime(at, base, humane.TimeOptions{IncludeSeconds: true})
			})

			context("just now", func() {
				it.BeforeEach(func() { at = ptr(base) })

				it("displays 0 seconds ago", func() {
					expect(result, t).To(Equal("0 seconds ago"))
				})
			})

			context("1 second ago", func() {
				it.BeforeEach(func() { at = ptr(base.Add(-1 * time.Second)) })

				it("displays 1 second ago, singular", func() {
					expect(result, t).To(Equal("1 second ago"))
				})
			})

			context("45 seconds ago", func() {
				it.BeforeEach(func() { at = ptr(base.Add(-45 * time.Second)) })

				it("displays 45 seconds ago", func() {
					expect(result, t).To(Equal("45 seconds ago"))
				})
			})

			context("45 seconds from now", func() {
				it.BeforeEach(func() { at = ptr(base.Add(45 * time.Second)) })

				it("displays in 45 seconds", func() {
					expect(result, t).To(Equal("in 45 seconds"))
				})
			})
		})

		describe("with Approximate: false", func() {
			var at *time.Time
			var result string
			it.JustBeforeEach(func() {
				result = humane.DistanceInTime(at, base, humane.TimeOptions{Approximate: humane.Bool(false)})
			})

			context("1 hour ago", func() {
				it.BeforeEach(func() { at = ptr(base.Add(-1 * time.Hour)) })

				it("displays the exact count, no about prefix", func() {
					expect(result, t).To(Equal("1 hour ago"))
				})
			})

			context("15 hours ago", func() {
				it.BeforeEach(func() { at = ptr(base.Add(-15 * time.Hour)) })

				it("displays 15 hours ago", func() {
					expect(result, t).To(Equal("15 hours ago"))
				})
			})
		})

		describe("nil handling", func() {
			var opts humane.TimeOptions
			var result string
			it.JustBeforeEach(func() { result = humane.DistanceInTime(nil, base, opts) })

			context("when at is nil and WhenNil is set", func() {
				it.BeforeEach(func() { opts = humane.TimeOptions{WhenNil: "an unknown time"} })

				it("returns WhenNil without formatting", func() {
					expect(result, t).To(Equal("an unknown time"))
				})
			})

			context("when at is nil and WhenNil is left unset", func() {
				it("returns an empty string", func() {
					expect(result, t).To(Equal(""))
				})
			})
		})

		// Boundary regression coverage for ActionView's distance_of_time_in_words bucket table (truncated at the "1 day" row); each context below sits on one cutoff second from that table.
		describe("at the approximate-distance bucket table boundaries", func() {
			context("with Approximate: false", func() {
				var at *time.Time
				var result string
				it.JustBeforeEach(func() {
					result = humane.DistanceInTime(at, base, humane.TimeOptions{Approximate: humane.Bool(false)})
				})

				context("29 seconds ago", func() {
					it.BeforeEach(func() { at = ptr(base.Add(-29 * time.Second)) })

					it("stays less than a minute", func() {
						expect(result, t).To(Equal("less than a minute ago"))
					})
				})

				context("30 seconds ago", func() {
					it.BeforeEach(func() { at = ptr(base.Add(-30 * time.Second)) })

					it("rounds up to 1 minute", func() {
						expect(result, t).To(Equal("1 minute ago"))
					})
				})

				context("89 seconds ago", func() {
					it.BeforeEach(func() { at = ptr(base.Add(-89 * time.Second)) })

					it("stays 1 minute", func() {
						expect(result, t).To(Equal("1 minute ago"))
					})
				})

				context("90 seconds ago", func() {
					it.BeforeEach(func() { at = ptr(base.Add(-90 * time.Second)) })

					it("rounds up to 2 minutes", func() {
						expect(result, t).To(Equal("2 minutes ago"))
					})
				})

				context("44 minutes 29 seconds ago", func() {
					it.BeforeEach(func() { at = ptr(base.Add(-(44*time.Minute + 29*time.Second))) })

					it("stays 44 minutes", func() {
						expect(result, t).To(Equal("44 minutes ago"))
					})
				})

				context("44 minutes 30 seconds ago", func() {
					it.BeforeEach(func() { at = ptr(base.Add(-(44*time.Minute + 30*time.Second))) })

					it("rounds up to 1 hour", func() {
						expect(result, t).To(Equal("1 hour ago"))
					})
				})

				context("89 minutes 29 seconds ago", func() {
					it.BeforeEach(func() { at = ptr(base.Add(-(89*time.Minute + 29*time.Second))) })

					it("stays 1 hour", func() {
						expect(result, t).To(Equal("1 hour ago"))
					})
				})

				context("89 minutes 30 seconds ago", func() {
					it.BeforeEach(func() { at = ptr(base.Add(-(89*time.Minute + 30*time.Second))) })

					it("rounds up to 2 hours", func() {
						expect(result, t).To(Equal("2 hours ago"))
					})
				})

				context("23 hours 59 minutes 29 seconds ago", func() {
					it.BeforeEach(func() { at = ptr(base.Add(-(23*time.Hour + 59*time.Minute + 29*time.Second))) })

					it("stays 24 hours", func() {
						expect(result, t).To(Equal("24 hours ago"))
					})
				})

				context("23 hours 59 minutes 30 seconds ago", func() {
					it.BeforeEach(func() { at = ptr(base.Add(-(23*time.Hour + 59*time.Minute + 30*time.Second))) })

					it("rounds up to 1 day", func() {
						expect(result, t).To(Equal("1 day ago"))
					})
				})
			})

			context("with no options", func() {
				var at *time.Time
				var result string
				it.JustBeforeEach(func() { result = humane.DistanceInTime(at, base, humane.TimeOptions{}) })

				context("44 minutes 29 seconds ago", func() {
					it.BeforeEach(func() { at = ptr(base.Add(-(44*time.Minute + 29*time.Second))) })

					it("has no about", func() {
						expect(result, t).To(Equal("44 minutes ago"))
					})
				})

				context("44 minutes 30 seconds ago", func() {
					it.BeforeEach(func() { at = ptr(base.Add(-(44*time.Minute + 30*time.Second))) })

					it("gains about, entering the hour bucket", func() {
						expect(result, t).To(Equal("about 1 hour ago"))
					})
				})

				context("23 hours 59 minutes 29 seconds ago", func() {
					it.BeforeEach(func() { at = ptr(base.Add(-(23*time.Hour + 59*time.Minute + 29*time.Second))) })

					it("keeps about", func() {
						expect(result, t).To(Equal("about 24 hours ago"))
					})
				})

				context("23 hours 59 minutes 30 seconds ago", func() {
					it.BeforeEach(func() { at = ptr(base.Add(-(23*time.Hour + 59*time.Minute + 30*time.Second))) })

					it("drops about, entering the day bucket", func() {
						expect(result, t).To(Equal("1 day ago"))
					})
				})
			})
		})
	})
}
