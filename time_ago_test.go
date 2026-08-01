package humane_test

import (
	"testing"
	"time"

	"github.com/sclevine/spec"
	. "github.com/woodie/expect"

	"github.com/woodie/humane"
)

// TimeAgo is a thin one-argument convenience over DistanceInTime, supplying
// time.Now() as relativeTo -- see DistanceInTime above for the exhaustive
// wording/bucket coverage this doesn't need to repeat. at is a plain
// time.Time (not *time.Time), so there's no nil-handling case here,
// unlike humane-ruby/humane-swift's time_ago/timeAgo.
func TestTimeAgo(t *testing.T) {
	spec.Run(t, "humane.TimeAgo", func(t *testing.T, context spec.G, it spec.S) {
		var when time.Time
		var result string
		it.JustBeforeEach(func() { result = humane.TimeAgo(when) })

		context("just now", func() {
			it.BeforeEach(func() { when = time.Now() })

			it("displays less than a minute ago", func() {
				expect(result, t).To(Equal("less than a minute ago"))
			})
		})

		context("3 minutes ago", func() {
			it.BeforeEach(func() { when = time.Now().Add(-3 * time.Minute) })

			it("forwards to DistanceInTime with time.Now() as relativeTo", func() {
				expect(result, t).To(Equal("3 minutes ago"))
			})
		})
	})
}
