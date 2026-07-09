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
		Context("with CollapseMinute (the default)", func() {
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

				It("displays less than a minute ago", func() {
					Expect(formatter.Format(when, base)).To(Equal("less than a minute ago"))
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

				It("displays in less than a minute", func() {
					Expect(formatter.Format(when, base)).To(Equal("in less than a minute"))
				})
			})

			Context("3 minutes from now", func() {
				BeforeEach(func() { when = base.Add(3 * time.Minute) })

				It("displays in 3 minutes", func() {
					Expect(formatter.Format(when, base)).To(Equal("in 3 minutes"))
				})
			})
		})

		Context("with CollapseMinute: false", func() {
			var (
				formatter humane.TimeFormatter
				when      time.Time
			)

			BeforeEach(func() {
				formatter = humane.TimeFormatter{}
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
	})
})
