package humane_test

import (
	"testing"

	"github.com/woodie/spec"
	. "github.com/woodie/expect"

	"github.com/woodie/humane"
)

func TestHumanSize(t *testing.T) {
	spec.Run(t, "humane.HumanSize", func(t *testing.T, context spec.G, it spec.S) {
		beforeEach, justBeforeEach := it.BeforeEach, it.JustBeforeEach

		// justBeforeEach runs the action under test (HumanSize) once, after
		// every beforeEach at every nesting level has set bytes -- the direct
		// replacement for a subject closure invoked explicitly in each it,
		// now that spec supports it natively.
		var bytes int64
		var result string
		justBeforeEach(func() { result = humane.HumanSize(bytes) })

		context("with 0 bytes", func() {
			beforeEach(func() { bytes = 0 })

			it("formats as Zero KB", func() {
				expect(result, t).To(Equal("Zero KB"))
			})
		})

		context("with 1 byte", func() {
			beforeEach(func() { bytes = 1 })

			it("spells out the singular unit", func() {
				expect(result, t).To(Equal("1 byte"))
			})
		})

		context("with a small byte count", func() {
			beforeEach(func() { bytes = 7 })

			it("spells out bytes", func() {
				expect(result, t).To(Equal("7 bytes"))
			})
		})

		context("with 999 bytes", func() {
			beforeEach(func() { bytes = 999 })

			it("stays in bytes", func() {
				expect(result, t).To(Equal("999 bytes"))
			})
		})

		context("with 79992 bytes", func() {
			beforeEach(func() { bytes = 79992 })

			it("formats as 80 KB", func() {
				expect(result, t).To(Equal("80 KB"))
			})
		})

		context("with a real file's byte count", func() {
			beforeEach(func() { bytes = 225935 })

			it("formats as 226 KB", func() {
				expect(result, t).To(Equal("226 KB"))
			})
		})

		context("with 500000 bytes", func() {
			beforeEach(func() { bytes = 500000 })

			it("formats as 500 KB", func() {
				expect(result, t).To(Equal("500 KB"))
			})
		})

		context("with a single-digit megabyte value", func() {
			beforeEach(func() { bytes = 1500000 })

			it("shows one decimal place, trailing zero trimmed", func() {
				expect(result, t).To(Equal("1.5 MB"))
			})
		})

		context("with a gigabyte-scale value", func() {
			beforeEach(func() { bytes = 5240000000 })

			it("shows two decimal places", func() {
				expect(result, t).To(Equal("5.24 GB"))
			})
		})

		context("with a value that lands on an exact unit", func() {
			beforeEach(func() { bytes = 2000000 })

			it("trims both trailing decimal digits", func() {
				expect(result, t).To(Equal("2 MB"))
			})
		})

		context("with a kilobyte value under 10", func() {
			beforeEach(func() { bytes = 1234 })

			it("rounds to whole kilobytes", func() {
				expect(result, t).To(Equal("1 KB"))
			})
		})

		context("with a kilobyte value that has a fraction", func() {
			beforeEach(func() { bytes = 2579 })

			it("rounds to the nearest kilobyte", func() {
				expect(result, t).To(Equal("3 KB"))
			})
		})

		context("with a megabyte value under 10", func() {
			beforeEach(func() { bytes = 4132383 })

			it("shows one decimal place, not two", func() {
				expect(result, t).To(Equal("4.1 MB"))
			})
		})

		context("with a megabyte value whose second decimal would round", func() {
			beforeEach(func() { bytes = 5243887 })

			it("drops the second decimal place", func() {
				expect(result, t).To(Equal("5.2 MB"))
			})
		})

		context("with a megabyte value at 100 or more", func() {
			beforeEach(func() { bytes = 123456789 })

			it("still shows one decimal place", func() {
				expect(result, t).To(Equal("123.5 MB"))
			})
		})

		context("with a gigabyte value at 10 or more", func() {
			beforeEach(func() { bytes = 12345678901 })

			it("still shows two decimal places", func() {
				expect(result, t).To(Equal("12.35 GB"))
			})
		})

		context("with a terabyte value", func() {
			beforeEach(func() { bytes = 1234567890123 })

			it("shows two decimal places", func() {
				expect(result, t).To(Equal("1.23 TB"))
			})
		})

		context("when rounding reaches 1000 of a unit", func() {
			context("just below the kilobyte limit", func() {
				beforeEach(func() { bytes = 999499 })

				it("stays in KB", func() {
					expect(result, t).To(Equal("999 KB"))
				})
			})

			context("with 999999 bytes", func() {
				beforeEach(func() { bytes = 999999 })

				it("moves up to MB", func() {
					expect(result, t).To(Equal("1 MB"))
				})
			})

			context("just below the megabyte limit", func() {
				beforeEach(func() { bytes = 999949999 })

				it("stays in MB", func() {
					expect(result, t).To(Equal("999.9 MB"))
				})
			})

			context("with 999960000 bytes", func() {
				beforeEach(func() { bytes = 999960000 })

				it("moves up to GB", func() {
					expect(result, t).To(Equal("1 GB"))
				})
			})

			context("with 999995000000 bytes", func() {
				beforeEach(func() { bytes = 999995000000 })

				it("moves up to TB", func() {
					expect(result, t).To(Equal("1 TB"))
				})
			})
		})
	})
}
