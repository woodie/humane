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
- **`TimeFormatter`**: asymmetric `"X ago"` / `"in X"` wording, matching
  `RelativeDateTimeFormatter`'s actual output exactly. `v0.1.0` shipped with
  symmetric `"X ago"` / `"X from now"` wording instead, documented at the time
  as "a deliberate departure, not an oversight" -- reverted in `v0.2.0` once it
  became clear that departure contradicted this library's own premise
  (matching what Swift/Finder-adjacent APIs actually do, the same bar
  `SizeFormatter` was held to via real-hardware comparison). No `"about"`
  prefix on the hour bucket (Go's `justincampbell/timeago`, still in `lambada`
  pre-integration, adds one; Swift's formatter doesn't either). `DateTimeStyle`/
  `.named` (`"yesterday"`, calendar-boundary-aware) isn't implemented -- genuine
  complexity, not trivial, and nothing downstream needs it yet.
- **`IncludeSeconds`** (bool, zero-value/default-intended `false`; renamed from
  `CollapseMinute` in `v0.3.0`, see `docs/releases/v0.3.0.md`): when `false`,
  renders anything under 60 seconds as `"less than a minute ago"`/`"in less
  than a minute"`. This collapsing doesn't exist in `RelativeDateTimeFormatter`
  at all -- zouk's own `ScanEntry.timeAgo` bolts a manual `< 30`-second clamp on
  top of the formatter for exactly this reason. Every real reference (Rails,
  Go's `timeago`, zouk's workaround) does this in practice, so there's no "pure
  Swift" behavior being overridden; the future-side wording follows the same
  asymmetric `"in X"` pattern as the counted buckets. Named and defaulted after
  ActionView's own `include_seconds`, which defaults the same way independently
  -- `humane-ruby` picked up the identical rename.
- **Go's zero-value gotcha -- resolved by the rename above**: `TimeFormatter{}`'s
  zero value used to have `CollapseMinute: false` (Go can't default a bare
  `bool` field to `true` from a struct literal) -- second-level granularity,
  the surprising case, not `NewTimeFormatter()`'s collapsed default. The
  `IncludeSeconds` rename inverted the polarity, so the zero value (`false`)
  *is* the recommended default now: `TimeFormatter{}` and `NewTimeFormatter()`
  are identical as of `v0.3.0`, and a test locks that in. `NewTimeFormatter()`
  is kept for API stability and parity with the other languages' constructors,
  not because the footgun still exists.

## Sandbox limitation

No Go toolchain in the Cowork sandbox (confirmed: no `go` binary) -- same
situation as `lambada`. Code changes here are made by inspection, then
**confirmed on real hardware**: woodie ran `go mod tidy`/`go test ./...`/
`ginkgo-fd` on their Mac each time, including after the switch to
Ginkgo/Gomega and the `lambada` integration below. Still true for anyone
picking this up fresh -- don't trust sandbox-only changes to `.go` files
until they've been run for real.

## Current state

Done: `SizeFormatter`, `TimeFormatter`, Ginkgo/Gomega tests, README,
`docs/COMMENTS.md` (long comments extracted per the convention in
zouk's `docs/COWORK.md`), a GitHub Actions `ci.yml`, and README badges.
Tagged and pushed as `v0.1.0`. Integrated into `lambada`'s `main.go`
(replacing `humanSize` and `justincampbell/timeago`-backed `timeAgo`),
released as `lambada` `2.2.0` -- confirmed via `go test ./...`, 44/44
passing. `humane-ruby` is the published Ruby sibling, integrated into
`scandalous` the same way, released as `scandalous` `2.2.0`.

`v0.2.0`: `TimeFormatter`'s future-side wording changed from symmetric
`"X from now"` to asymmetric `"in X"`, matching `RelativeDateTimeFormatter`
exactly -- see "Design decisions" above. Breaking change to the string
output; `lambada` and `scandalous` (and their test suites, and zouk's own
`ScanEntry` if it ever compares wording) need a follow-up pass once this
is tagged and published, since they're currently locked to the old
`"X from now"` wording.

`v0.3.0`: `TimeFormatter.CollapseMinute` (recommended `true` via
`NewTimeFormatter()`) renamed to `IncludeSeconds` (zero value `false`) -- an
exact polarity inversion, so the recommended default behavior is unchanged.
As a side effect, this retires the long-standing Go zero-value gotcha:
`TimeFormatter{}` and `NewTimeFormatter()` are now identical, since
`IncludeSeconds`' zero value is itself the recommended default -- see
"Design decisions" above and `docs/releases/v0.3.0.md`. `humane-ruby` picked
up the identical rename in its own `v0.3.0` (already committed, not yet
tagged). Not yet confirmed via `go test ./...` -- no Go toolchain in this
sandbox, made by inspection only; needs a real run on woodie's Mac before
tagging.

## Next up

1. Confirm `v0.3.0` via `go test ./...` on real hardware, then tag and push.
2. Propagate both the `v0.2.0` wording change and `v0.3.0` rename into
   `lambada` and `scandalous`.
3. Decide whether `humane-swift`'s `approximate` option (ActionView-inspired
   "about"/"in about" prefixing on hour-plus buckets) gets backported here
   too -- see `humane-swift/docs/COWORK.md` "Next up".
4. `SizeFormatter` has no `AllowedUnits`/`CountStyle` (Finder's style is the
   only one anything downstream needs today), and `TimeFormatter` has no
   `.named` style (`"yesterday"`, calendar-boundary-aware) -- both left out
   deliberately per "Design decisions" above, not gaps to fill without a
   real need.
5. `humane-swift`'s real-hardware testing found `ByteCountFormatter`'s actual
   output diverges from this package's hand-rolled 2-significant-digit math
   in a few cases (zero bytes, byte-scale labels, some GB-scale precision) --
   see `humane-swift/docs/COWORK.md` "Current state" for specifics. Worth
   deciding whether to correct `SizeFormatter` toward exact parity or
   document the gap as accepted.
