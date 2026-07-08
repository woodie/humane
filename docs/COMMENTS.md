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
Formats one time relative to another the way Finder-adjacent tools do:
symmetric "X ago" / "X from now" phrasing, with no "about" prefix on the
hour bucket -- Swift's RelativeDateTimeFormatter has no such prefix
either, and Go's justincampbell/timeago (which does add one) is exactly
what this package replaces.

The zero value has CollapseMinute set to false (second-level granularity
for anything under a minute); use NewTimeFormatter for the collapsed
default described below. Go structs can't default a bool field to true
from a bare literal, so the constructor is the recommended way to get it
-- TimeFormatter{CollapseMinute: true} works identically if you'd rather
not use it.

### `TimeFormatter.CollapseMinute`
Renders any duration under 60 seconds as "less than a minute ago" /
"less than a minute from now" instead of counting seconds. Rails'
distance_of_time_in_words, Go's justincampbell/timeago, and zouk's own
RelativeDateTimeFormatter wrapper all do this in practice -- Swift's
formatter has no such bucket natively, so there's no "pure" behavior
being overridden here, just a convenience every real reference already
reaches for.

### `NewTimeFormatter`
Returns a TimeFormatter with CollapseMinute enabled, the recommended
default.

### `TimeFormatter.Format`
Returns t relative to relativeTo as a human-readable string.

    f := NewTimeFormatter()
    f.Format(t, t)                    == "less than a minute ago"
    f.Format(t, t.Add(3*time.Minute))  == "3 minutes ago"
    f.Format(t, t.Add(-3*time.Minute)) == "3 minutes from now"
    f.Format(t, t.Add(15*time.Hour))   == "15 hours ago"
    f.Format(t, t.Add(30*time.Hour))   == "1 day ago"
