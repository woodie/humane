# Comments

Rationale, history, and design notes that used to live as multi-line
comments in the source. Organized by file, then by the type, property, or
function each note is attached to. The source itself now carries at most
one short line at any given spot -- anything longer that would previously
have been a doc comment lives here instead.

## size.go

### Package `humane`
Formats file sizes and relative dates the way macOS Finder does, as
package-level functions (`HumanSize`, `TimeAgo`) rather than configurable
formatter types. Through `v0.5.0` this package mirrored Foundation's
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

### `HumanSize`
Formats bytes the way Finder does. Bakes in three corrections found by
comparing this package's original 2-significant-figure port against real
`ByteCountFormatter` output (via `humane-swift`'s real-hardware testing --
see that repo's `docs/COWORK.md`, "Current state"):

- `0` reads `"Zero KB"`, not `"0 B"` -- a hardcoded special case;
  `ByteCountFormatter` doesn't run its usual rounding logic for zero.
- Values under 1000 spell out `"byte"`/`"bytes"` (`"1 byte"`, `"7 bytes"`)
  rather than using a `"B"` label.
- Everything else is rounded to 3 significant figures (not 2) via
  `formatSignificant`, then unit-labeled. The old 2-significant-figure rule
  (1 decimal below 10, none at or above) undercounted precision for
  values under 10: `5,240,000,000` bytes is `"5.24 GB"` on real hardware,
  not `"5.2 GB"`.

The 3-significant-figure rule was chosen over a narrower "just fix the
GB case" patch because it's the only single rule found that reproduces
every known fixture at once -- including the two cross-checked against
real hardware before this change (`225_935` -> `"226 KB"`, `500_000` ->
`"500 KB"`) and the existing `1,500,000` -> `"1.5 MB"` fixture, alongside
the new GB finding. That said, it's still an inference from a small
fixture set, not something confirmed across every magnitude -- see
`docs/COWORK.md`'s `v0.9.0` entry for what still needs a real
`ByteCountFormatter` comparison.

    HumanSize(0)          == "Zero KB"
    HumanSize(1)          == "1 byte"
    HumanSize(79992)      == "80 KB"
    HumanSize(225935)     == "226 KB"
    HumanSize(1500000)    == "1.5 MB"
    HumanSize(5240000000) == "5.24 GB"

### `formatSignificant`
Rounds `value` to `sigFigs` significant figures, then trims trailing
fractional zeros (and the decimal point itself, if nothing remains after
it) -- this trimming is what keeps `"1.5 MB"` from becoming `"1.50 MB"`
while still letting `"5.24 GB"` keep both of its non-zero decimal digits.
`magnitude` (digits before the decimal point) uses `floor(log10(value)) +
1` rather than `ceil(log10(value))` specifically to avoid a boundary bug
at exact powers of 10 (`ceil(log10(10))` and `ceil(log10(10.0))` can both
round to `1` instead of the correct `2`-digit magnitude, depending on
floating-point rounding of the logarithm itself).

## time.go

### `TimeOptions` (struct)
Replaces the `TimeFormatter` struct/`NewTimeFormatter()` constructor pair
from `v0.5.0` and earlier. Passed as a trailing variadic argument to
`TimeAgo` specifically so omitting it entirely still gets the recommended
defaults (`TimeAgo(at, relativeTo)` with no third argument) without the
zero-value ambiguity a single non-variadic `TimeOptions{}` parameter would
have -- see `TimeOptions.Approximate` below for why that ambiguity would
otherwise matter here.

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
Added in `v0.9.0` alongside `TimeAgo` accepting `at *time.Time` instead of
`at time.Time`. Motivated by `zouk`'s `ScanEntry.timeAgo(relativeTo:)`,
which used to guard a possibly-unparsable timestamp itself
(`guard let downloadedAt else { return nil }`) and hand the caller a
`String?` that still needed its own `?? "an unknown time"` fallback one
layer up in `ScanGridView` -- two guard points for one final string.
`TimeAgo` now takes the optional directly and a caller-supplied fallback
string, collapsing both guard points into one call. The fallback text
stays app-specific (an empty default, not a hardcoded "unknown time" or
similar baked into this package) -- consistent with keeping
ActionView-flavored vocabulary opt-in and configurable rather than
assumed, the same principle `Approximate`/`IncludeSeconds` already follow.

### `TimeAgo`
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

    TimeAgo(at, at)                              == "less than a minute ago"
    TimeAgo(at, at.Add(45*time.Second))           == "1 minute ago"
    TimeAgo(at, at.Add(3*time.Minute))            == "3 minutes ago"
    TimeAgo(at, at.Add(-3*time.Minute))           == "in 3 minutes"
    TimeAgo(at, at.Add(15*time.Hour))             == "about 15 hours ago"
    TimeAgo(at, at.Add(30*time.Hour))             == "1 day ago" // no "about" -- ActionView's table has none on the day bucket
    TimeAgo(nil, now, TimeOptions{WhenNil: "an unknown time"}) == "an unknown time"
