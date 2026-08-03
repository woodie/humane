# Picking up humane in a new Cowork session

Cross-project conventions (git locks, sandbox toolchain gaps, pushing, comments, code
style) are in `~/workspace/woodie/docs/COWORK.md`.

## What this is

A small Go module formatting file sizes and relative dates the way macOS Finder
does, matching Swift's `ByteCountFormatter`/`RelativeDateTimeFormatter` output plus
an opt-in layer of ActionView-style vocabulary on top. Package-level functions, no
instantiation: `HumanSize(bytes)`, `DistanceInTime(at, relativeTo, opts...)`, and
`TimeAgo(at, opts...)` (a one-argument convenience wrapping `DistanceInTime` with
`time.Now()`). Siblings: `humane-ruby`, `humane-swift`, `humane-kotlin` -- same
algorithm, same wording, each in its own repo/registry.

## Design decisions worth knowing

- **`TimeOptions.Approximate` is `*bool`, not `bool`.** Go can't default a bare
  `bool` struct field to `true` from a zero-value literal, so a plain `bool` would
  make `TimeOptions{}` silently mean "off" the moment a caller writes any other
  explicit field -- the same zero-value gotcha `IncludeSeconds` used to have under
  its old name (`CollapseMinute`). `*bool` lets `TimeOptions{}`'s zero value still
  mean "use the default (`true`)".
- **`WhenNil` collapses a guard-then-fallback into one call** -- `TimeAgo`/
  `DistanceInTime` take `at *time.Time` directly rather than requiring the caller to
  nil-check first, motivated by `zouk`'s Swift `ScanEntry.timeAgo` doing exactly that
  guard by hand before this existed.
- **Defaults match Foundation exactly**; ActionView's `Approximate`/`IncludeSeconds`
  vocabulary is opt-in, layered on top, never a silent behavior change to the
  baseline.
- **No `replace` directive needed for `woodie/spec`** the way `gorderly` avoids one
  -- `humane` is a library, nothing ever runs `go install` against it, so the
  `go install`-rejects-`replace` problem that forced `gorderly` off the fork doesn't
  apply here. `go.mod` currently pins `woodie/spec v0.2.0` via `replace` for
  `BeforeEach`/`AfterEach`/`JustBeforeEach`.

## Testing

`sclevine/spec` (via the `woodie/spec` fork) + `github.com/woodie/expect`, default
output format `-fs` (matching the rest of the family's wording/spec structure).
`go test -v ./... | gorderly -fd`, or `make test`/`make check`.

## Sandbox limitation

No Go toolchain here -- changes are written by inspection, verified via `go mod
tidy`/`go test ./...`/`make check` on the user's own Mac.

## Current status

`v0.9.5`. Fully migrated off Ginkgo/Gomega onto `spec`+`expect` (matches `gorderly`/
`lambada`). Adopted by `lambada`, `scandalous` (via `humane-ruby`), and `zouk` (via
`humane-swift`).

## Open items

**`HumanSize` rounding verification** is tracked on `humane-swift` (the real-hardware
reference implementation), not here -- see that repo's own `docs/COWORK.md`.
`humane-ruby` issue #1 (ActionView's full bucket table past the "1 day" row) is a
deliberate scope boundary, not a gap -- see that repo's issue tracker, not a TODO
here.
