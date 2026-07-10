# Comments

Rationale, history, and design notes that used to live as multi-line
comments in the source. Organized by file, then by the type, property, or
function each note is attached to. The source itself now carries at most
one short line at any given spot -- anything longer that would previously
have been a doc comment lives here instead.

## size.go

### Package `humane`
Formats file sizes and relative dates the way macOS Finder does, modeled
on Swift's ByteCountFormatter and RelativeDateTimeFormatter: small,
configurable formatter types with a single Format method, rather than a
bare function.

### `sizeUnits`
Finder's capitalized, 1000-based unit labels -- not the SI-correct
lowercase "kB", and not 1024-based "KiB"/"MiB". No published Go library
ships this exact combination; see the README for why.

### `SizeFormatter` (struct)
Formats byte counts the way Finder does: 1000-based math, capitalized
unit labels, rounded to 2 significant digits. The zero value is ready to
use -- there's no configuration yet, since Finder's is the only style
this package currently needs.

### `SizeFormatter.Format`
Returns bytes as a Finder-style human-readable string.

    SizeFormatter{}.Format(79992)   == "80 KB"
    SizeFormatter{}.Format(225935)  == "226 KB"
    SizeFormatter{}.Format(1500000) == "1.5 MB"

## time.go

### `TimeFormatter` (struct)
Formats one time relative to another the way RelativeDateTimeFormatter
does: asymmetric "X ago" / "in X" phrasing (matched exactly, not
"X from now" -- an earlier symmetric wording was found to be an
unforced departure from the very API this package is modeled on). No
"about" prefix on the hour bucket by default -- Swift's
RelativeDateTimeFormatter has no such prefix either, and Go's
justincampbell/timeago (which does add one) is exactly what this package
replaces -- but see Approximate below for an explicit opt-in.

Renamed CollapseMinute to IncludeSeconds in v0.3.0 (see
docs/releases/v0.3.0.md) -- an exact polarity inversion, which happens to
retire the zero-value footgun this section used to warn about: the zero
value now has IncludeSeconds set to false, and false is the collapsed
(recommended) behavior under the new name, so TimeFormatter{} and
NewTimeFormatter() are now identical. Under the old name, the zero value
(CollapseMinute: false) was the surprising second-level-granularity case;
that asymmetry is gone.

### `TimeFormatter.IncludeSeconds`
When false (the zero value and the default), renders any duration under
60 seconds as "less than a minute ago" / "in less than a minute" instead
of counting seconds. Rails' distance_of_time_in_words, Go's
justincampbell/timeago, and zouk's own RelativeDateTimeFormatter wrapper
all do this in practice -- Swift's formatter has no such bucket natively,
so there's no "pure" behavior being overridden here, just a convenience
every real reference already reaches for. The future phrasing follows
the same asymmetric "in X" pattern as the counted buckets below, not a
symmetric "X from now". Named and defaulted after ActionView's own
include_seconds, which defaults the same way.

### `NewTimeFormatter`
Returns a TimeFormatter with the recommended default -- kept for API
stability and parity with the other two languages' constructors, even
though it's now equivalent to TimeFormatter{} (see above).

### `TimeFormatter.Approximate`
Added in v0.4.0 (see docs/releases/v0.4.0.md). When true, prefixes
"about"/"in about" onto any bucket of an hour or larger -- matching
ActionView's distance_of_time_in_words past that same boundary. Sub-hour
buckets are untouched either way. Defaults to false, matching Foundation's
raw output.

Ported from humane-swift's identically-named, identically-defaulted
option (v0.1.0), for contexts that render once and can't refresh (a web
response) where a precise-looking "15 hours ago" overstates the value's
actual precision.

Format already builds a bare quantity phrase (text) before wrapping it in
"X ago"/"in X", so prefixing "about " onto text first composes correctly
in both directions with no string surgery -- the same shape humane-ruby's
#string uses. Swift's TimeFormatter has to post-process
RelativeDateTimeFormatter's already-complete phrase instead, since
Foundation hands back the whole sentence at once.

### `TimeFormatter.Format`
Returns t relative to relativeTo as a human-readable string.

    f := NewTimeFormatter()
    f.Format(t, t)                    == "less than a minute ago"
    f.Format(t, t.Add(3*time.Minute))  == "3 minutes ago"
    f.Format(t, t.Add(-3*time.Minute)) == "in 3 minutes"
    f.Format(t, t.Add(15*time.Hour))   == "15 hours ago"
    f.Format(t, t.Add(30*time.Hour))   == "1 day ago"

    a := TimeFormatter{Approximate: true}
    a.Format(t, t.Add(15*time.Hour))   == "about 15 hours ago"
    a.Format(t, t.Add(30*time.Hour))   == "about 1 day ago"
    a.Format(t, t.Add(-3*time.Hour))   == "in about 3 hours"
