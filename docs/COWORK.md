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
- **`woodie/spec` is a plain `require`, no `replace` directive.** `woodie/spec`
  v0.3.0 renamed its module declaration to its own path (`woodie/spec#3`), so
  `go.mod` pins `github.com/woodie/spec v0.3.0` directly for
  `BeforeEach`/`AfterEach`/`JustBeforeEach`.

## Testing

`woodie/spec` + `github.com/woodie/expect`, default output format `-fs`
(matching the rest of the family's wording/spec structure). `go test -v
./... | gorderly -fd`, or `make test`/`make check`.

**Deliberate exception: hook methods are aliased here, not called
qualified.** Every other Go repo in this account (`gorderly`, `lambada`,
`expect` itself) calls `it.BeforeEach`/`it.JustBeforeEach` qualified --
`spec`'s own `docs/COWORK.md` explains why the fuller destructuring
alias (`context, before, after := describe, it.Before, it.After`) got
reversed once three hook names made that line cluttered. `humane`
intentionally goes back to that fuller style, as a real worked example
of it: `describe, beforeEach, justBeforeEach := context, it.BeforeEach,
it.JustBeforeEach` (or just `beforeEach, justBeforeEach := ...` in files
that don't call `describe`). Works cleanly here specifically because
this suite never uses `AfterEach` -- two aliased hook names, not three,
stays readable where the account-wide reversal was reacting to three.
See `config_test.go`'s comment for the same note at the point of use.
Not a mistake to "fix" into qualified calls to match sibling repos --
this repo is deliberately the exception, on purpose, for comparison.

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
