package humane_test

import (
	"testing"

	"github.com/sclevine/spec"
	. "github.com/woodie/expect"

	"github.com/woodie/humane"
)

func TestHumanSize(t *testing.T) {
	spec.Run(t, "HumanSize", func(t *testing.T, context spec.G, it spec.S) {
		// JustBeforeEach runs the action under test (HumanSize) once, after
		// every BeforeEach at every nesting level has set bytes -- the direct
		// replacement for a subject closure invoked explicitly in each it,
		// now that spec supports it natively.
		var bytes int64
		var result string
		it.JustBeforeEach(func() { result = humane.HumanSize(bytes) })

		context("with 0 bytes", func() {
			it.BeforeEach(func() { bytes = 0 })

			it("formats as Zero KB, matching ByteCountFormatter's own wording", func() {
				expect(result, t).To(Equal("Zero KB"))
			})
		})

		context("with 1 byte", func() {
			it.BeforeEach(func() { bytes = 1 })

			it("spells out the singular unit", func() {
				expect(result, t).To(Equal("1 byte"))
			})
		})

		context("with a small byte count", func() {
			it.BeforeEach(func() { bytes = 7 })

			it("spells out bytes rather than using a B label", func() {
				expect(result, t).To(Equal("7 bytes"))
			})
		})

		context("with 999 bytes", func() {
			it.BeforeEach(func() { bytes = 999 })

			it("stays in bytes, just under the 1 KB threshold", func() {
				expect(result, t).To(Equal("999 bytes"))
			})
		})

		context("with the shared 79992-byte fixture used by lambada scandalous", func() {
			it.BeforeEach(func() { bytes = 79992 })

			it("formats as 80 KB", func() {
				expect(result, t).To(Equal("80 KB"))
			})
		})

		context("with a real file's byte count", func() {
			it.BeforeEach(func() { bytes = 225935 })

			it("matches Finder's reported size", func() {
				expect(result, t).To(Equal("226 KB"))
			})
		})

		context("with zouk's ByteCountFormatter(.file) fixture", func() {
			it.BeforeEach(func() { bytes = 500000 })

			it("matches its output", func() {
				expect(result, t).To(Equal("500 KB"))
			})
		})

		context("with a single-digit megabyte value", func() {
			it.BeforeEach(func() { bytes = 1500000 })

			it("shows one decimal place, trailing zero trimmed", func() {
				expect(result, t).To(Equal("1.5 MB"))
			})
		})

		context("with a gigabyte-scale value", func() {
			it.BeforeEach(func() { bytes = 5240000000 })

			it("keeps 2 decimal places at 3 significant figures (not truncated to 1)", func() {
				expect(result, t).To(Equal("5.24 GB"))
			})
		})

		context("with a value that lands on an exact unit", func() {
			it.BeforeEach(func() { bytes = 2000000 })

			it("trims both trailing decimal digits", func() {
				expect(result, t).To(Equal("2 MB"))
			})
		})
	})
}
