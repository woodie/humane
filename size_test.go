package humane_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/woodie/humane"
)

var _ = Describe("HumanSize", func() {
	Context("with 0 bytes", func() {
		It("formats as Zero KB, matching ByteCountFormatter's own wording", func() {
			Expect(humane.HumanSize(0)).To(Equal("Zero KB"))
		})
	})

	Context("with 1 byte", func() {
		It("spells out the singular unit", func() {
			Expect(humane.HumanSize(1)).To(Equal("1 byte"))
		})
	})

	Context("with a small byte count", func() {
		It("spells out bytes rather than using a B label", func() {
			Expect(humane.HumanSize(7)).To(Equal("7 bytes"))
		})
	})

	Context("with 999 bytes", func() {
		It("stays in bytes, just under the 1 KB threshold", func() {
			Expect(humane.HumanSize(999)).To(Equal("999 bytes"))
		})
	})

	Context("with the shared 79992-byte fixture used by lambada/scandalous", func() {
		It("formats as 80 KB", func() {
			Expect(humane.HumanSize(79992)).To(Equal("80 KB"))
		})
	})

	Context("with a real file's byte count", func() {
		It("matches Finder's reported size", func() {
			Expect(humane.HumanSize(225935)).To(Equal("226 KB"))
		})
	})

	Context("with zouk's ByteCountFormatter(.file) fixture", func() {
		It("matches its output", func() {
			Expect(humane.HumanSize(500000)).To(Equal("500 KB"))
		})
	})

	Context("with a single-digit megabyte value", func() {
		It("shows one decimal place, trailing zero trimmed", func() {
			Expect(humane.HumanSize(1500000)).To(Equal("1.5 MB"))
		})
	})

	Context("with a gigabyte-scale value", func() {
		It("keeps 2 decimal places at 3 significant figures (not truncated to 1)", func() {
			Expect(humane.HumanSize(5240000000)).To(Equal("5.24 GB"))
		})
	})

	Context("with a value that lands on an exact unit", func() {
		It("trims both trailing decimal digits", func() {
			Expect(humane.HumanSize(2000000)).To(Equal("2 MB"))
		})
	})
})
