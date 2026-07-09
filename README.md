# humane

[![go.mod version](https://img.shields.io/github/go-mod/go-version/woodie/humane)](https://github.com/woodie/humane)
[![CI](https://github.com/woodie/humane/actions/workflows/ci.yml/badge.svg)](https://github.com/woodie/humane/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/woodie/humane.svg)](https://github.com/woodie/humane/releases/latest)
[![License](https://img.shields.io/github/license/woodie/humane.svg)](LICENSE)

Getting human-readable file sizes with 1000-based math
(as the Mac Finder displays) and relative times worded the way Swift's
`RelativeDateTimeFormatter` does turned out to be a real challenge to get
both right and simple. The `humane` library exists so a Go application can share
consistent size and time formatting with a Swift application, instead of
reaching for a library whose output doesn't match Swift's or that's
complicated to drop in.

```go
import "github.com/woodie/humane"

humane.SizeFormatter{}.Format(225935) // "226 KB"

timeFormatter := humane.NewTimeFormatter()
timeFormatter.Format(time.Now().Add(-3*time.Minute), time.Now()) // "3 minutes ago"
```

Corresponding functions in Swift will have consistent output.

```swift
import Foundation

ByteCountFormatter.string(fromByteCount: Int64(225935), countStyle: .file) // "226 KB"

let formatter = RelativeDateTimeFormatter(); formatter.unitsStyle = .full
formatter.localizedString(for: time, relativeTo: now) // "3 minutes ago"
```

## Install

```
go get github.com/woodie/humane
```

## Scope

Finder's `.file` byte-count style, and a numeric (non-calendar-aware)
relative time style -- that's the whole surface area today.
`AllowedUnits`/alternate `CountStyle`s and `.named` style (`"yesterday"`,
calendar-boundary-aware) aren't implemented -- contributions welcome.
