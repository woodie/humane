package humane_test

import (
	"testing"

	. "github.com/woodie/expect"
)

// Allow all tests in this package to use lowercase expect()
func expect[T any](got T, t testing.TB) Expectation[T] { return Expect(got, t) }

// Improve readability with structural functions and lifecycle hooks.
// Parameter is named context spec.G (not describe) -- only add
// `describe := context` if a file actually calls describe(...), e.g.
// distance_in_time_test.go. human_size_test.go/time_ago_test.go only
// ever call context(...) and skip the alias entirely.
// Deliberate exception to the account default: this repo also aliases
// it.BeforeEach/it.JustBeforeEach to lowercase locals (beforeEach,
// justBeforeEach), the fuller destructuring style used before hook
// names were reversed to qualified calls elsewhere -- kept here on
// purpose as a worked example, since two hook names (no AfterEach in
// this suite) stays readable where three would clutter.
// https://gist.github.com/woodie/35ee3fc2bea01b775b95b3fe5e950a05#file-example-go-L3
