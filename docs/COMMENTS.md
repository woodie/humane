# Comments

Rationale, history, and design notes that used to live as multi-line
comments in the source. Organized by file, then by the type, property, or
function each note is attached to. The source itself now carries at most
one short line at any given spot -- anything longer that would previously
have been a doc comment lives here instead.

## size.go

### Package `humane`
Formats file sizes and relative dates the way macOS Finder does, as
package-level functions (`HumanSize`, `DistanceInTime`, `TimeAgo`) rather
than configurable formatter types. Through `v0.5.0` this package mirrored Foundation's
formatter-type shape (`SizeFormatter{}.Format(...)`, `NewTimeFormatter()`)
on purpose, matching `ByteCountFormatter`/`RelativeDateTimeFormatter`'s own
API surface. `v0.9.0` drops that shape: once configuration moved to a
per-call options argument (see `time.go`), a formatter instance held no
state at all, so instantiating one first was ceremony left over from a
design this package no longer follows. See `docs/COWORK.md`'s `v0.9.0`
entry for the full rationale -- consistency with ActionView's own bare
helper-function shape (`number_to_human_size`, `distance_of_time_in_words`)
matters more now than mirroring Foundation's API surface, since the actual
goal is dropping into Ruby/Go HTML templates as simply as ActionView does.

### `sizeUnits`
Finder's capitalized, 1000-based unit labels -- not the SI-correct
lowercase "kB", and not 1024-based "KiB"/"MiB". No published Go library
ships this exact combination; see the README for why. Starts at `KB`
(index 0) -- byte-scale values are spelled out (`"7 bytes"`) rather than
using a unit label at all, so `sizeUnits` never needs a `"B"` entry.

### `decimalPlaces`
Finder's fixed decimal places per unit, index-aligned with `sizeUnits`:
`KB` 0, `MB` 1, `GB` and above 2. This is `ByteCountFormatter`'s default
("adaptive") behavior, confirmed against real Finder output in issue #3.

### `HumanSize`
Formats bytes the way Finder does:

- `0` reads `"Zero KB"`, not `"0 B"` -- a hardcoded special case;
  `ByteCountFormatter` doesn't run its usual rounding logic for zero.
- Values under 1000 spell out `"byte"`/`"bytes"` (`"1 byte"`, `"7 bytes"`)
  rather than using a `"B"` label.
- Everything else is scaled to its unit, rounded half up to that unit's
  `decimalPlaces`, then has trailing fractional zeros trimmed (`"2 MB"`,
  not `"2.0 MB"`).
- If rounding reaches 1000 of a unit, the value moves up to the next unit
  and is re-rounded there: `999_999` is `"1 MB"`, not `"1000 KB"`, and
  `999_960_000` is `"1 GB"`, not `"1000 MB"`.

Earlier versions rounded everything to 3 significant figures instead. That
rule happened to match every fixture the suite had at the time, but real
Finder output showed it was wrong for most sizes under 10 MB (`4,132,383`
bytes is `"4.1 MB"`, not `"4.13 MB"`) and above 100 MB (`123,456,789` is
`"123.5 MB"`, not `"123 MB"`). See issue #3 for the full Finder comparison.

    HumanSize(0)           == "Zero KB"
    HumanSize(1)           == "1 byte"
    HumanSize(2579)        == "3 KB"
    HumanSize(225935)      == "226 KB"
    HumanSize(999999)      == "1 MB"
    HumanSize(4132383)     == "4.1 MB"
    HumanSize(123456789)   == "123.5 MB"
    HumanSize(5240000000)  == "5.24 GB"
    HumanSize(12345678901) == "12.35 GB"

### `roundForUnit`
Uses `math.Round` (half away from zero) on the scaled value rather than
letting `strconv.FormatFloat` do the rounding, since `FormatFloat` rounds
exact ties to even. `1500` bytes is `"2 KB"` in Finder.

### `trimZeros`
Formats an already-rounded value to its unit's decimal places, then trims
trailing fractional zeros (and the decimal point itself, if nothing remains
after it).

## time.go

### `TimeOptions` (struct)
Replaces the `TimeFormatter` struct/`NewTimeFormatter()` constructor pair
from `v0.5.0` and earlier. Passed as a trailing variadic argument to
`DistanceInTime` (and, since `v0.9.3`, `TimeAgo`) specifically so omitting
it entirely still gets the recommended defaults (`DistanceInTime(at,
relativeTo)` with no third argument) without the zero-value ambiguity a
single non-variadic `TimeOptions{}` parameter would have -- see
`TimeOptions.Approximate` below for why that ambiguity would otherwise
matter here.

### `TimeOptions.Approximate`
A `*bool`, not a `bool`. `v0.9.0` flips the recommended default for
`Approximate` from `false` to `true`, matching ActionView's own
`distance_of_time_in_words` (which has no toggle for this at all -- it's
always on past the hour boundary) and, in practice, matching what every
real consumer (`lambada`, `scandalous`, `zouk`) already passed explicitly
under the old API. If `Approximate` were a plain `bool`, `TimeOptions{}`'s
zero value would be `false` -- silently the opposite of the new
recommended default the moment any caller writes an explicit `TimeOptions{
IncludeSeconds: true}` without also repeating `Approximate: true`. This is
the exact class of bug `IncludeSeconds` itself used to have under its old
name `CollapseMinute` (see `docs/releases/v0.3.0.md`) -- fixed there by
inverting the field's polarity so `false` became the safe zero value.
`Approximate` can't be fixed the same way (there's no natural opposite
name/polarity that reads well), so it's a pointer instead: `nil` means
"use the default (`true`)", and `Bool(false)` (see below) opts out
explicitly.

### `Bool`
A one-line pointer-of-a-literal helper, needed only because Go has no
inline syntax for taking the address of a literal (`&false` is a syntax
error). Exists solely so `TimeOptions.Approximate` can be set to an
explicit `false` without a caller declaring their own local variable
first.

### `TimeOptions.WhenNil`
Added in `v0.9.0` alongside `TimeAgo` (renamed `DistanceInTime` in `v0.9.3`)
accepting `at *time.Time` instead of `at time.Time`. Motivated by `zouk`'s
`ScanEntry.timeAgo(relativeTo:)`, which used to guard a possibly-unparsable
timestamp itself (`guard let downloadedAt else { return nil }`) and hand
the caller a `String?` that still needed its own `?? "an unknown time"`
fallback one layer up in `ScanGridView` -- two guard points for one final
string. `DistanceInTime`/`TimeAgo` now take the optional directly and a
caller-supplied fallback string, collapsing both guard points into one
call. The fallback text stays app-specific (an empty default, not a
hardcoded "unknown time" or similar baked into this package) -- consistent
with keeping ActionView-flavored vocabulary opt-in and configurable rather
than assumed, the same principle `Approximate`/`IncludeSeconds` already
follow.

### `DistanceInTime` (renamed from `TimeAgo` in `v0.9.3`)
Formats one time relative to another the way ActionView's
`distance_of_time_in_words` does for wording, but direction-aware like
`RelativeDateTimeFormatter` -- `"X ago"` for the past, `"in X"` for the
future, chosen automatically from `relativeTo.Sub(*at)`'s sign rather than
requiring the caller to know which one applies ahead of time (which is
what ActionView itself requires, and what `distance_of_time_in_words`'s
own `.abs` collapses future distances into a past-tense string as a known,
unfixed bug -- see `humane-ruby`'s original `docs/COWORK.md` "Why this
exists" section).

Buckets are chosen from `distanceInMinutes` (seconds/60, rounded once via
`int(d.Minutes()+0.5)`), not by re-dividing raw seconds independently per
unit -- the old per-unit approach let rounding carry across a bucket
boundary on its own (`59:59:59` used to round to `"60 minutes ago"`
instead of `"1 hour ago"`). Computing `distanceInMinutes` once and
switching on *that* is exactly how ActionView's own
`distance_of_time_in_words` works, and is what produces its specific,
non-obvious cutoffs: the "about 1 hour" bucket starts at 44 minutes 30
seconds (not 60:00), and "about 2 hours" starts at 89:30, not 90:00.

    DistanceInTime(at, at)                              == "less than a minute ago"
    DistanceInTime(at, at.Add(45*time.Second))           == "1 minute ago"
    DistanceInTime(at, at.Add(3*time.Minute))            == "3 minutes ago"
    DistanceInTime(at, at.Add(-3*time.Minute))           == "in 3 minutes"
    DistanceInTime(at, at.Add(15*time.Hour))             == "about 15 hours ago"
    DistanceInTime(at, at.Add(30*time.Hour))             == "1 day ago" // no "about" -- ActionView's table has none on the day bucket
    DistanceInTime(nil, now, TimeOptions{WhenNil: "an unknown time"}) == "an unknown time"

### `TimeAgo` (new meaning in `v0.9.3`)
A one-argument convenience over `DistanceInTime`, supplying `time.Now()` as
`relativeTo` -- see `docs/COWORK.md`'s `v0.9.3` entry for why (the
ActionView `distance_of_time_in_words`/`time_ago_in_words` naming pair this
mirrors). Thin passthrough, no bucketing logic of its own; the `TimeAgo`
`Describe` block in `time_test.go` covers it with three cases rather than
re-testing everything `DistanceInTime` already does.
