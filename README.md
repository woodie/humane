# humane

[![go.mod version](https://img.shields.io/github/go-mod/go-version/woodie/humane)](https://github.com/woodie/humane)
[![CI](https://github.com/woodie/humane/actions/workflows/ci.yml/badge.svg)](https://github.com/woodie/humane/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/woodie/humane.svg)](https://github.com/woodie/humane/releases/latest)
[![License](https://img.shields.io/github/license/woodie/humane.svg)](LICENSE)

File sizes that read exactly the way macOS Finder shows them (`"4.1 MB"`,
`"226 KB"`), and relative dates in Rails ActionView's wording (`"about 3 hours
ago"`, `"in 3 minutes"`), for Go -- as simple to drop into an HTML template as
ActionView's own helpers.

```go
import "github.com/woodie/humane"

humane.HumanSize(225935) // "226 KB"

mtime := time.Now().Add(-180 * time.Second)
humane.TimeAgo(mtime) // "3 minutes ago" -- relative to the real clock

now := time.Now()
humane.DistanceInTime(&mtime, now) // "3 minutes ago" -- explicit relativeTo, for tests
```

## Install

```
go get github.com/woodie/humane
```

## File sizes

`HumanSize` matches Foundation's `ByteCountFormatter` (file style), which is what
macOS Finder shows, and is checked against real Finder output:

- 1000-based units with capitalized labels: `KB`, `MB`, `GB`, `TB`, and up.
- A fixed number of decimal places per unit -- none for KB, one for MB, two
  for GB and above -- with trailing zeros dropped (`"2 MB"`, not `"2.0 MB"`).
- Rounds half up, and moves up a unit when rounding reaches 1000:
  999,999 bytes is `"1 MB"`, not `"1000 KB"`.
- Zero bytes is `"Zero KB"`, and anything under 1000 bytes is spelled out
  (`"1 byte"`, `"7 bytes"`).

```go
humane.HumanSize(2579)       // "3 KB"
humane.HumanSize(4132383)    // "4.1 MB"
humane.HumanSize(999999)     // "1 MB"
humane.HumanSize(5240000000) // "5.24 GB"
```

## `DistanceInTime` and `TimeAgo`

Two entry points, same naming split as ActionView's own
`distance_of_time_in_words`/`time_ago_in_words`:

- **`DistanceInTime(at *time.Time, relativeTo, ...)`** -- the explicit,
  fully-tested core. Takes both times, plus a possibly-nil `at`, so it's
  what tests should call.
- **`TimeAgo(at time.Time, ...)`** -- a one-argument convenience for the
  common "drop into a view" case, where a timestamp is always in hand.
  Supplies `time.Now()` as `relativeTo` internally.

The wording follows ActionView's `distance_of_time_in_words`, with one
difference: direction is picked for you. A past time reads `"3 minutes ago"`
and a future one `"in 3 minutes"`, where ActionView leaves that to the caller.

Both share the same options and recommended defaults, already matching
ActionView's own `distance_of_time_in_words` defaults -- pass none at all and
you get them for free:

```go
humane.DistanceInTime(at, relativeTo) // Approximate: true, IncludeSeconds: false
humane.TimeAgo(at)                    // same defaults, relativeTo is time.Now()
```

- **`Approximate`** (`*bool`, default `true`): prefixes `"about"`/`"in about"`
  on the hour-scale buckets (1 hour, and 2..24 hours), matching ActionView's
  `distance_of_time_in_words` wording for those buckets exactly (down to its
  44:30/89:30 rounding cutoffs), through the "1 day" bucket. This is a
  `*bool`, not a `bool` -- see "Why `Approximate` is a pointer" below.
- **`IncludeSeconds`** (`bool`, default `false`): under 30 seconds, collapses
  to `"less than a minute ago"`/`"in less than a minute"` instead of an exact
  second count. Matches ActionView's `include_seconds` default.
- **`WhenNil`** (`string`, default `""`): if `at` is `nil`, `DistanceInTime`
  returns this string without formatting -- for a scan, download, or other
  record that doesn't have a timestamp yet. `TimeAgo` takes `at` by value
  (never nil); use `DistanceInTime` directly for a possibly-missing timestamp.

```go
humane.DistanceInTime(at, relativeTo, humane.TimeOptions{Approximate: humane.Bool(false)})
// "15 hours ago", not "about 15 hours ago"

humane.DistanceInTime(nil, time.Now(), humane.TimeOptions{WhenNil: "an unknown time"})
// "an unknown time"
```

### Why `Approximate` is a pointer

`TimeOptions{}`'s zero value needs to mean "use the recommended defaults,"
and Go can't default a bare `bool` field to `true` from a struct literal --
the zero value of `bool` is always `false`. (This is the same footgun that
`IncludeSeconds` used to have under its old name, `CollapseMinute`; see
`docs/COMMENTS.md`.) Rather than silently defaulting `Approximate` to
`false` whenever you write an explicit `TimeOptions{...}` for some other
field, `Approximate` is a `*bool`: `nil` means "use the default (`true`)",
and `humane.Bool(false)` opts out.

## Scope

Finder's byte-count style, and a numeric (non-calendar-aware) relative time
style through the "1 day" bucket -- that's the whole surface area today.
Alternate size units/styles and a `.named` style (`"yesterday"`,
calendar-boundary-aware) aren't implemented -- contributions welcome.

## Development

```
golangci-lint run
go test -v ./... | gorderly -fd
```

or, if you don't have [`gorderly`](https://github.com/woodie/gorderly)
installed:

```
go vet ./...
go test ./...
```
