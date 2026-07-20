package humane_test

import (
	"testing"

	. "github.com/woodie/expect"
)

// expect is the lowercase call-site alias recommended in expect's own
// README ("Lowercase call sites") -- a one-line generic pass-through
// declared once per test package, since Go's capitalize-to-export rule
// only applies across the package boundary. Keeps every call site here
// reading lowercase alongside describe/context/it/before instead of
// standing out as the one capitalized word in the block, with zero loss of
// compile-time type inference. Shared by size_test.go and time_test.go.
func expect[T any](got T, t testing.TB) Expectation[T] { return Expect(got, t) }

// Improve readability with vars set for structural functions and lifecycle hooks
/*
func TestCalculator(t *testing.T) {
	spec.Run(t, "Calculator", func(t *testing.T, describe spec.G, it spec.S) {
		context, before, after := describe, it.Before, it.After // HERE

		var calculator *Calculator
		before(func() { calculator = NewCalculator() })
		after(func() { calculator.Clear() })

		describe("#add", func() {
			context("with positive numbers", func() {
				it("returns the correct sum", func() {
					expect(calculator.Add(2, 3), t).To(Equal(5))
				})
			})
		})
	})
} */
