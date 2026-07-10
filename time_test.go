package humane_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/woodie/humane"
)

var _ = Describe("TimeFormatter", func() {
	base := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)

	Describe("Format", func() {
		Context("with IncludeSeconds: false (the default)", func() {
			var (
				formatter humane.TimeFormatter
				when      time.Time
			)

			BeforeEach(func() {
				formatter = humane.NewTimeFormatter()
			})

			Context("just now", func() {
				BeforeEach(func() { when = base })

				It("displays less than a minute ago", func() {
					Expect(formatter.Format(when, base)).To(Equal("less than a minute ago"))
				})
			})

			Context("45 seconds ago", func() {
				BeforeEach(func() { when = base.Add(-45 * time.Second) })

				It("rounds up to 1 minute ago (past the 30-second cutoff)", func() {
					Expect(formatter.Format(when, base)).To(Equal("1 minute ago"))
				})
			})

			Context("1 minute ago", func() {
				BeforeEach(func() { when = base.Add(-1 * time.Minute) })

				It("displays 1 minute ago, singular", func() {
					Expect(formatter.Format(when, base)).To(Equal("1 minute ago"))
				})
			})

			Context("3 minutes ago", func() {
				BeforeEach(func() { when = base.Add(-3 * time.Minute) })

				It("displays 3 minutes ago", func() {
					Expect(formatter.Format(when, base)).To(Equal("3 minutes ago"))
				})
			})

			Context("1 hour ago", func() {
				BeforeEach(func() { when = base.Add(-1 * time.Hour) })

				It("displays 1 hour ago, singular", func() {
					Expect(formatter.Format(when, base)).To(Equal("1 hour ago"))
				})
			})

			Context("15 hours ago", func() {
				BeforeEach(func() { when = base.Add(-15 * time.Hour) })

				It("displays 15 hours ago", func() {
					Expect(formatter.Format(when, base)).To(Equal("15 hours ago"))
				})
			})

			Context("30 hours ago", func() {
				BeforeEach(func() { when = base.Add(-30 * time.Hour) })

				It("rolls up to 1 day ago", func() {
					Expect(formatter.Format(when, base)).To(Equal("1 day ago"))
				})
			})

			Context("3 days ago", func() {
				BeforeEach(func() { when = base.Add(-3 * 24 * time.Hour) })

				It("displays 3 days ago", func() {
					Expect(formatter.Format(when, base)).To(Equal("3 days ago"))
				})
			})

			Context("45 seconds from now", func() {
				BeforeEach(func() { when = base.Add(45 * time.Second) })

				It("rounds up to in 1 minute (past the 30-second cutoff)", func() {
					Expect(formatter.Format(when, base)).To(Equal("in 1 minute"))
				})
			})

			Context("3 minutes from now", func() {
				BeforeEach(func() { when = base.Add(3 * time.Minute) })

				It("displays in 3 minutes", func() {
					Expect(formatter.Format(when, base)).To(Equal("in 3 minutes"))
				})
			})
		})

		Context("with IncludeSeconds: true", func() {
			var (
				formatter humane.TimeFormatter
				when      time.Time
			)

			BeforeEach(func() {
				formatter = humane.TimeFormatter{IncludeSeconds: true}
			})

			Context("just now", func() {
				BeforeEach(func() { when = base })

				It("displays 0 seconds ago", func() {
					Expect(formatter.Format(when, base)).To(Equal("0 seconds ago"))
				})
			})

			Context("1 second ago", func() {
				BeforeEach(func() { when = base.Add(-1 * time.Second) })

				It("displays 1 second ago, singular", func() {
					Expect(formatter.Format(when, base)).To(Equal("1 second ago"))
				})
			})

			Context("45 seconds ago", func() {
				BeforeEach(func() { when = base.Add(-45 * time.Second) })

				It("displays 45 seconds ago", func() {
					Expect(formatter.Format(when, base)).To(Equal("45 seconds ago"))
				})
			})

			Context("45 seconds from now", func() {
				BeforeEach(func() { when = base.Add(45 * time.Second) })

				It("displays in 45 seconds", func() {
					Expect(formatter.Format(when, base)).To(Equal("in 45 seconds"))
				})
			})
		})

		Context("with Approximate: true", func() {
			var (
				formatter humane.TimeFormatter
				when      time.Time
			)

			BeforeEach(func() {
				formatter = humane.TimeFormatter{Approximate: true}
			})

			Context("59 minutes ago", func() {
				BeforeEach(func() { when = base.Add(-59 * time.Minute) })

				It("prefixes about (59 minutes falls in the 45..89-minute 'about 1 hour' bucket)", func() {
					Expect(formatter.Format(when, base)).To(Equal("about 1 hour ago"))
				})
			})

			Context("exactly 1 hour ago", func() {
				BeforeEach(func() { when = base.Add(-1 * time.Hour) })

				It("displays about 1 hour ago", func() {
					Expect(formatter.Format(when, base)).To(Equal("about 1 hour ago"))
				})
			})

			Context("15 hours ago", func() {
				BeforeEach(func() { when = base.Add(-15 * time.Hour) })

				It("displays about 15 hours ago", func() {
					Expect(formatter.Format(when, base)).To(Equal("about 15 hours ago"))
				})
			})

			Context("30 hours ago", func() {
				BeforeEach(func() { when = base.Add(-30 * time.Hour) })

				It("does not prefix about on the day bucket (ActionView's table has no 'about 1 day')", func() {
					Expect(formatter.Format(when, base)).To(Equal("1 day ago"))
				})
			})

			Context("3 minutes from now", func() {
				BeforeEach(func() { when = base.Add(3 * time.Minute) })

				It("leaves sub-hour buckets untouched", func() {
					Expect(formatter.Format(when, base)).To(Equal("in 3 minutes"))
				})
			})

			Context("3 hours from now", func() {
				BeforeEach(func() { when = base.Add(3 * time.Hour) })

				It("displays in about 3 hours", func() {
					Expect(formatter.Format(when, base)).To(Equal("in about 3 hours"))
				})
			})
		})

		// Boundary regression coverage for the ActionView distance_of_time_in_words bucket
		// table this approximate-distance behavior ports, truncated at the "1 day" row
		// since month/year buckets are out of scope. Each pair straddles a cutoff second
		// from that table to lock in exactly where the wording flips.
		Context("at the approximate-distance bucket table boundaries", func() {
			Context("without Approximate", func() {
				formatter := humane.NewTimeFormatter()

				It("29s stays less than a minute, 30s rounds up to 1 minute", func() {
					Expect(formatter.Format(base.Add(-29*time.Second), base)).To(Equal("less than a minute ago"))
					Expect(formatter.Format(base.Add(-30*time.Second), base)).To(Equal("1 minute ago"))
				})

				It("89s stays 1 minute, 90s rounds up to 2 minutes", func() {
					Expect(formatter.Format(base.Add(-89*time.Second), base)).To(Equal("1 minute ago"))
					Expect(formatter.Format(base.Add(-90*time.Second), base)).To(Equal("2 minutes ago"))
				})

				It("44:29 stays 44 minutes, 44:30 rounds up to 1 hour", func() {
					Expect(formatter.Format(base.Add(-(44*time.Minute+29*time.Second)), base)).To(Equal("44 minutes ago"))
					Expect(formatter.Format(base.Add(-(44*time.Minute+30*time.Second)), base)).To(Equal("1 hour ago"))
				})

				It("89:29 stays 1 hour, 89:30 rounds up to 2 hours", func() {
					Expect(formatter.Format(base.Add(-(89*time.Minute+29*time.Second)), base)).To(Equal("1 hour ago"))
					Expect(formatter.Format(base.Add(-(89*time.Minute+30*time.Second)), base)).To(Equal("2 hours ago"))
				})

				It("23:59:29 stays 24 hours, 23:59:30 rounds up to 1 day", func() {
					Expect(formatter.Format(base.Add(-(23*time.Hour+59*time.Minute+29*time.Second)), base)).To(Equal("24 hours ago"))
					Expect(formatter.Format(base.Add(-(23*time.Hour+59*time.Minute+30*time.Second)), base)).To(Equal("1 day ago"))
				})
			})

			Context("with Approximate: true", func() {
				formatter := humane.TimeFormatter{Approximate: true}

				It("44:29 has no about, 44:30 gains about (entering the hour bucket)", func() {
					Expect(formatter.Format(base.Add(-(44*time.Minute+29*time.Second)), base)).To(Equal("44 minutes ago"))
					Expect(formatter.Format(base.Add(-(44*time.Minute+30*time.Second)), base)).To(Equal("about 1 hour ago"))
				})

				It("23:59:29 keeps about, 23:59:30 drops about (entering the day bucket)", func() {
					Expect(formatter.Format(base.Add(-(23*time.Hour+59*time.Minute+29*time.Second)), base)).To(Equal("about 24 hours ago"))
					Expect(formatter.Format(base.Add(-(23*time.Hour+59*time.Minute+30*time.Second)), base)).To(Equal("1 day ago"))
				})
			})
		})

		Describe("TimeFormatter{} vs NewTimeFormatter()", func() {
			It("produce identical output, now that IncludeSeconds' zero value is the recommended default", func() {
				Expect(humane.TimeFormatter{}.Format(base, base)).To(Equal(humane.NewTimeFormatter().Format(base, base)))
			})
		})
	})
})
