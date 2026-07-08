# Picking up humane in a new Cowork session

Context for whoever opens this repo cold, with none of the prior conversation history.
Cross-project conventions (git locks, sandbox toolchain gaps, pushing, comments, code
style) are in `~/workspace/woodie/docs/COWORK.md`.

## What this is

A small Go module for formatting file sizes and relative dates the way macOS Finder
does, modeled on Swift's `ByteCountFormatter` and `RelativeDateTimeFormatter`:
configurable formatter *types* with a `Format` method, rather than a bare helper
function (the shape most Go/Ruby humanize libraries actually use).
[`humane-ruby`](https://github.com/woodie/humane-ruby) is the Ruby sibling -- same
algorithm, same wording, separate repo since Go module versioning and RubyGems
versioning don't share a tag namespace cleanly.

## Why this exists

Extracted out of `lambada` and `scandalous` after a multi-step saga fixing their
file-size formatting (see lambada's `docs/COWORK.md`, the "humanSize/timeAgo went
through N shapes" history, steps 6-8, for the full blow-by-blow):

1. lambada showed `"80 kB"` (go-humanize, SI/1000-based, lowercase), scandalous
   showed `"78.1 KB"` (Rails' `number_to_human_size`, 1024-based, capitalized).
   Fixed lambada to match scandalous's label -- wrong move, see next.
2. Real-world testing (comparing a live file's size across lambada-web,
   scandalous-web, zouk, and actual Finder) revealed scandalous's `"78.1 KB"` was
   *also* wrong: Rails' `number_to_human_size` is 1024-based despite the `KB`
   label, and Finder is 1000-based under that same label. Only zouk (Swift's
   `ByteCountFormatter(.file)`) was right the whole time, for free, via the OS.
3. No published Go or Ruby library ships 1000-based math under capitalized
   `KB`/`MB` labels -- the SI-correct ones (`go-humanize`, `docker/go-units`) use
   lowercase `kB`; the ones that capitalize it pair it with 1024-based math. Both
   lambada and scandalous ended up hand-rolling the same fix independently.
4. Since the same fix had to be written twice, `humane`/`humane-ruby` exist to
   write it once. Picked up `TimeFormatter` while at it, replacing lambada's
   `justincampbell/timeago` and scandalous's
   `ActionView::Helpers::DateHelper#time_ago_in_words` -- both had their own
   wording quirks (the "about" prefix, Ruby's future-date bug) a shared, tested
   implementation sidesteps.

## Naming

`cocoa`, `aqua`, `finder`, and `cupertino` were all checked and are already taken
on RubyGems.org (unrelated old gems). `humane` was open -- double meaning:
human-readable formatting, and a nod to Apple's Human Interface Guidelines, which
is the actual design lineage here. The Ruby sibling is `humane-ruby` rather than
`humane-gem`: "-ruby" names the language (matching the
`google-cloud-go`/`google-cloud-ruby` convention), "-gem" would've just said
"this is a gem," true of every RubyGem and not distinguishing information.

## Design decisions

- **`SizeFormatter`**: zero-config, zero value ready to use. No
  `AllowedUnits`/`CountStyle` -- Apple's real options aren't trivial to replicate
  faithfully across every magnitude bucket, and there's exactly one style
  (Finder's) anything here actually needs. Algorithm is a straight port of
  `go-humanize`'s `Bytes()` rounding (2 significant digits: one decimal when the
  integer part is a single digit, none once it hits two), just with 1000-based
  math (unchanged from `go-humanize`) and capitalized labels (`go-humanize` uses
  lowercase `kB`).
- **`TimeFormatter`**: symmetric `"X ago"` / `"X from now"` wording -- a
  deliberate departure from `RelativeDateTimeFormatter`'s actual asymmetric
  output (`"X ago"` / `"in X"`), not an oversight. No `"about"` prefix on the
  hour bucket (Go's `justincampbell/timeago`, still in `lambada`
  pre-integration, adds one; Swift's formatter doesn't either). `DateTimeStyle`/
  `.named` (`"yesterday"`, calendar-boundary-aware) isn't implemented -- genuine
  complexity, not trivial, and nothing downstream needs it yet.
- **`CollapseMinute`** (bool, default-intended `true`): renders anything under
  60 seconds as `"less than a minute ago"`/`"...from now"`. Doesn't exist in
  `RelativeDateTimeFormatter` at all -- zouk's own `ScanEntry.timeAgo` bolts a
  manual `< 30`-second clamp on top of the formatter for exactly this reason.
  Every real reference (Rails, Go's `timeago`, zouk's workaround) does this in
  practice, so there's no "pure Swift" behavior being overridden.
- **Go's zero-value gotcha**: `TimeFormatter{}`'s zero value has
  `CollapseMinute: false` (Go can't default a bare `bool` field to `true` from a
  struct literal) -- second-level granularity, not the collapsed default.
  `NewTimeFormatter()` returns the intended default (`CollapseMinute: true`).
  Flagged explicitly in the README and the struct doc comment so it doesn't read
  as a bug later.

## Sandbox limitation

No Go toolchain in the Cowork sandbox (confirmed: no `go` binary) -- same
situation as `lambada`. `size.go`/`time.go`/their tests were written by
inspection and hand-traced against the test cases by arithmetic, not run.
**Not yet confirmed**: ask whoever's continuing this to run `go test ./...` on a
machine with Go and report back before trusting this is bug-free.

## Current state

`SizeFormatter` and `TimeFormatter` are implemented with tests (`size_test.go`,
`time_test.go`), not yet confirmed passing (see Sandbox limitation).
`humane-ruby` hasn't been started. Once both are done and confirmed, `lambada`'s
`main.go` and `scandalous`'s `web.rb` both need to swap their hand-rolled
formatters for these packages (dropping `justincampbell/timeago` from
`lambada`'s `go.mod` and `ActionView::Helpers::DateHelper` from `scandalous`'s
`web.rb` in the process).

## Next up

- Confirm `go test ./...` passes on real Go.
- Build `humane-ruby`: same `SizeFormatter`/`TimeFormatter` shape, same
  algorithm, same wording, RSpec tests mirroring these Go tests.
- Integrate into `lambada` (replace `humanSize`/`timeAgo` in `main.go`) and
  `scandalous` (replace `human_size`/`time_ago_in_words` in `web.rb`), updating
  each repo's test fixtures for the wording changes (e.g., `"about 15 hours
  ago"` -> `"15 hours ago"`).
- Publish: tag `humane` as a Go module (works automatically once pushed with a
  semver tag), `gem push` `humane-ruby` to RubyGems.
