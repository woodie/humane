package humane_test

import (
	"testing"

	"github.com/sclevine/spec"
	. "github.com/woodie/expect"

	"github.com/woodie/humane"
)

func TestHumanSize(t *testing.T) {
	spec.RunAliased(t, "HumanSize", func(t *testing.T, _, context spec.Describe, it spec.S, _, _ func(func())) {
		context("with 0 bytes", func() {
			subject := humane.HumanSize(0)

			it("formats as Zero KB, matching ByteCountFormatter's own wording", func() {
				Expect(t, subject).To(Equal("Zero KB"))
			})
		})

		context("with 1 byte", func() {
			subject := humane.HumanSize(1)

			it("spells out the singular unit", func() {
				Expect(t, subject).To(Equal("1 byte"))
			})
		})

		context("with a small byte count", func() {
			subject := humane.HumanSize(7)

			it("spells out bytes rather than using a B label", func() {
				Expect(t, subject).To(Equal("7 bytes"))
			})
		})

		context("with 999 bytes", func() {
			subject := humane.HumanSize(999)

			it("stays in bytes, just under the 1 KB threshold", func() {
				Expect(t, subject).To(Equal("999 bytes"))
			})
		})

		context("with the shared 79992-byte fixture used by lambada scandalous", func() {
			subject := humane.HumanSize(79992)

			it("formats as 80 KB", func() {
				Expect(t, subject).To(Equal("80 KB"))
			})
		})

		context("with a real file's byte count", func() {
			subject := humane.HumanSize(225935)

			it("matches Finder's reported size", func() {
				Expect(t, subject).To(Equal("226 KB"))
			})
		})

		context("with zouk's ByteCountFormatter(.file) fixture", func() {
			subject := humane.HumanSize(500000)

			it("matches its output", func() {
				Expect(t, subject).To(Equal("500 KB"))
			})
		})

		context("with a single-digit megabyte value", func() {
			subject := humane.HumanSize(1500000)

			it("shows one decimal place, trailing zero trimmed", func() {
				Expect(t, subject).To(Equal("1.5 MB"))
			})
		})

		context("with a gigabyte-scale value", func() {
			subject := humane.HumanSize(5240000000)

			it("keeps 2 decimal places at 3 significant figures (not truncated to 1)", func() {
				Expect(t, subject).To(Equal("5.24 GB"))
			})
		})

		context("with a value that lands on an exact unit", func() {
			subject := humane.HumanSize(2000000)

			it("trims both trailing decimal digits", func() {
				Expect(t, subject).To(Equal("2 MB"))
			})
		})
	})
}
