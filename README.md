# humane
Swift's file sizes and relative dates for Go

Finder-accurate file sizes and relative dates for Go, modeled on Swift's [`ByteCountFormatter`](https://developer.apple.com/documentation/foundation/bytecountformatter) and [`RelativeDateTimeFormatter`](https://developer.apple.com/documentation/foundation/relativedatetimeformatter) -- not literal ports (both are closed-source, and `TimeFormatter`'s wording is a deliberate departure), but the same idea: a small, configurable formatter object instead of a bare function.

## Install

```
go get github.com/woodie/humane
```

## Usage

```go
import "github.com/woodie/humane"

sizeFormatter := humane.SizeFormatter{}
sizeFormatter.Format(225935) // "226 KB" -- 1000-based math, capitalized
                              // units, matching Finder, not the SI-correct
                              // lowercase "kB" or 1024-based "KiB"

timeFormatter := humane.NewTimeFormatter() // CollapseMinute: true
timeFormatter.Format(scannedAt, time.Now())
// "3 minutes ago" / "3 minutes from now" / "less than a minute ago"
```

`SizeFormatter{}`'s zero value is ready to use. `TimeFormatter{}`'s zero
value has `CollapseMinute: false` (second-level granularity under a
minute) since Go can't default a bool field to `true` from a bare
struct literal -- use `NewTimeFormatter()`, or set the field explicitly,
for the collapsed default described above.

## Scope

Only what lambada and scandalous actually need today: Finder's `.file`
byte-count style, and a numeric (non-calendar-aware) relative time
style. `ByteCountFormatter`'s `allowedUnits`/alternate `countStyle`s and
`RelativeDateTimeFormatter`'s `.named` style (`"yesterday"`, calendar-
boundary-aware) aren't implemented -- contributions welcome.
