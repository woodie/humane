package humane_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/woodie/humane"
)

var _ = Describe("SizeFormatter", func() {
	var formatter humane.SizeFormatter

	BeforeEach(func() {
		formatter = humane.SizeFormatter{}
	})

	Describe("Format", func() {
		Context("with 0 bytes", func() {
			It("formats as 0 B", func() {
				Expect(formatter.Format(0)).To(Equal("0 B"))
			})
		})

		Context("with a small byte count", func() {
			It("formats with no rounding", func() {
				Expect(formatter.Format(7)).To(Equal("7 B"))
			})
		})

		Context("with 999 bytes", func() {
			It("stays in bytes, just under the 1 KB threshold", func() {
				Expect(formatter.Format(999)).To(Equal("999 B"))
			})
		})

		Context("with the shared 79992-byte fixture used by lambada/scandalous", func() {
			It("formats as 80 KB", func() {
				Expect(formatter.Format(79992)).To(Equal("80 KB"))
			})
		})

		Context("with a real file's byte count", func() {
			It("matches Finder's reported size", func() {
				Expect(formatter.Format(225935)).To(Equal("226 KB"))
			})
		})

		Context("with zouk's ByteCountFormatter(.file) fixture", func() {
			It("matches its output", func() {
				Expect(formatter.Format(500000)).To(Equal("500 KB"))
			})
		})

		Context("with a single-digit megabyte value", func() {
			It("shows one decimal place", func() {
				Expect(formatter.Format(1500000)).To(Equal("1.5 MB"))
			})
		})

		Context("with a gigabyte-scale value", func() {
			It("rounds to 2 significant digits", func() {
				Expect(formatter.Format(5240000000)).To(Equal("5.2 GB"))
			})
		})
	})
})
