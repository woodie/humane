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
  prefix on the hour bucket by default (Go's `justincampbell/timeago`, still
  in `lambada` pre-integration, adds one; Swift's formatter doesn't either) --
  see `Approximate` below for the `v0.4.0` opt-in. `DateTimeStyle`/`.named`
  (`"yesterday"`, calendar-boundary-aware) isn't implemented -- genuine
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
- **`Approximate`** (bool, zero-value/default `false`; added in `v0.4.0`, see
  `docs/releases/v0.4.0.md`): prefixes `"about"`/`"in about"` onto any bucket
  of an hour or larger, matching ActionView's `distance_of_time_in_words`
  past that same boundary. Sub-hour buckets are untouched either way. Ported
  from `humane-swift`'s identically-named, identically-defaulted option
  (`v0.1.0`), for contexts that render once and can't refresh (a web
  response) where a precise-looking `"15 hours ago"` overstates the value's
  actual precision. `Format` already builds a bare quantity phrase (`text`)
  before wrapping it in `"X ago"`/`"in X"`, so prefixing `"about "` onto
  `text` composes correctly in both directions with no string surgery --
  `humane-ruby`'s `#string` has the same shape. Swift's `TimeFormatter` has
  to post-process `RelativeDateTimeFormatter`'s already-complete phrase
  instead, since Foundation hands back the whole sentence at once.

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
up the identical rename in its own `v0.3.0`, tagged, pushed, and published
to RubyGems. Confirmed for real on woodie's Mac: `go test ./...` (`ok`) and
`ginkgo-fd` -- 23/23, including the new `TimeFormatter{} vs
NewTimeFormatter()` equivalence spec. Tagged and released as `v0.3.0`:
https://github.com/woodie/humane/releases/tag/v0.3.0.

`v0.4.0`: `TimeFormatter` gains `Approximate` (zero value `false`) -- see
"Design decisions" above and `docs/releases/v0.4.0.md`. Additive, not
breaking. Ported from `humane-swift`'s `v0.1.0` option, following
`humane-ruby`'s identical `v0.4.0` addition (tagged and published to
RubyGems, wired into `scandalous`'s `time_ago` helper, released as
`scandalous` `2.5.0`). New Ginkgo context added to `time_test.go` covering
the hour boundary, a 15-hour and 30-hour past case, and a 3-hour future
case. Confirmed via real `go test ./...`/`ginkgo-fd` on woodie's Mac, then
tagged, pushed, and released: https://github.com/woodie/humane/releases/tag/v0.4.0.

Wired into `lambada`'s `timeFormatter` in the same window (`Approximate:
true`, mirroring `scandalous`), `go.mod`/`go.sum` bumped to `v0.4.0` via a
real `go get`/`go mod tidy`, confirmed via `go test ./...`/`ginkgo-fd`
(45/45) plus `npm run check` (JS lint/tests, `golangci-lint`), released as
`lambada` `2.5.0`, and deployed to the Pi -- confirmed live (`"about 14
hours ago"`).

Also this session: README's Swift code sample (raw `ByteCountFormatter`/
`RelativeDateTimeFormatter` calls) now links to `humane-swift` directly
instead of only showing bare Foundation calls, and gained the "Beyond
Foundation's defaults" section documenting `IncludeSeconds`/`Approximate`
that this README had never had (unlike `humane-ruby`'s and
`humane-swift`'s READMEs, which already did).

`v0.5.0`: `Format` reworked to match the ActionView
`distance_of_time_in_words` bucket table quoted in `humane-ruby` issue #1
exactly, through the "1 day" row -- see `humane-ruby`'s own `docs/COWORK.md`
for the full writeup and rationale, ported here identically. `IncludeSeconds:
false`'s collapse cutoff moved 60s -> 30s, `Approximate` narrowed from "any
bucket >= 1 hour" to exactly the hour-scale buckets, and `Format` now buckets
off `distanceInMinutes` rather than raw seconds re-divided per unit (fixes a
latent bug where `59:59:59` rounded to "60 minutes ago"). New boundary-pair
specs added covering all six table cutoffs. Confirmed for real on woodie's
Mac: `go test ./...` (`ok`) and `ginkgo-fd` -- 36/36 passing, alongside
`humane-ruby`'s identical change (35/35) in the same session. Tagged and
released: https://github.com/woodie/humane/releases/tag/v0.5.0.

`humane-ruby` issue #1 (https://github.com/woodie/humane-ruby/issues/1,
"Provide ActionView compatibility mode") is the source of the bucket table
above and remains open: it quotes ActionView's full table, including the
2..29-day and month/year buckets past the "1 day" row this package
implements. No further work is scheduled against it beyond what `v0.5.0`
already ported -- see "Next up" below.

## `v0.9.0`: a full API rethink, informed by three real consumers

Through `v0.5.0` this package was built as a Go port of `ByteCountFormatter`/
`RelativeDateTimeFormatter` -- configurable formatter types mirroring
Foundation's own API shape. With three real consumers now shipped
(`lambada`, `scandalous`, `zouk`) and this package the sole user of its own
API (woodie), the goal shifted: instead of *looking like* Foundation, feel
as simple to drop into a Go or Ruby HTML template as ActionView's own
helpers do, while still behaving consistently across all three languages.
That's a different design bar than "port Foundation faithfully," and this
release rebuilds the public API around it rather than iterating on the old
shape further. Breaking changes throughout -- acceptable since nothing
outside this session's own `lambada`/`scandalous`/`zouk` follow-up depends
on the old API yet.

**`SizeFormatter{}.Format(...)`/`NewTimeFormatter()` -> `HumanSize(...)`/
`TimeAgo(...)`, package-level functions, no instantiation.** Once
configuration moved to a per-call options argument (below), a formatter
struct held no state between calls -- instantiating one first was ceremony
left over from mirroring Foundation, not something this design still needs.
`humane-ruby` and `humane-swift` picked up the equivalent change in the
same session (class methods, static methods respectively) -- each language
uses its own idiomatic casing for the shared concept (`HumanSize`/`TimeAgo`
here, `human_size`/`time_ago` in Ruby, `humanSize`/`timeAgo` in Swift)
rather than one literal spelling forced across all three. See
`docs/COMMENTS.md`.

**`HumanSize`'s rounding corrected toward 3 significant figures.** See
`docs/COMMENTS.md` for the full derivation -- bakes in `humane-swift`'s
real-hardware findings (`"Zero KB"` for zero, spelled-out `"bytes"` below
1000) plus a rounding-rule change (3 significant figures, trailing zeros
trimmed) that reproduces every known fixture, old and new, with one rule.
Resolves item 2 from the old "Next up" list below.

**`TimeAgo`'s `Approximate` default flips `false` -> `true`.** Matches
ActionView's own `distance_of_time_in_words` (which has no toggle for this
at all -- always on past the hour boundary), and, checked against real
code, matches what every current consumer already passes explicitly
(`lambada`'s `main.go`, `scandalous`'s `web.rb`, `zouk`'s `ScanEntry.swift`
all set it `true` under the old API). Zero behavior change for any shipped
consumer; removes required boilerplate at every call site instead. Needed
a new mechanism to avoid a **second Go zero-value gotcha**: `TimeOptions`
replaces the old `TimeFormatter` struct, and `Approximate` is a `*bool`
rather than a `bool` specifically so `TimeOptions{}`'s zero value still
means "use the default (`true`)" instead of silently meaning `false` the
moment a caller writes an explicit `TimeOptions{IncludeSeconds: true}` --
see `docs/COMMENTS.md` for why this is the same bug class `IncludeSeconds`
itself used to have under its old name, `CollapseMinute`.

**`TimeAgo` takes `at *time.Time` and a `WhenNil` option.** Motivated by
`zouk`'s `ScanEntry.timeAgo(relativeTo:)`, which guarded a possibly-missing
timestamp itself and handed the caller a `String?` needing its own `??`
fallback one layer up -- two guard points for one string. `TimeAgo` now
takes the optional directly; `TimeOptions.WhenNil` supplies the fallback
text in the same call, collapsing both layers into one. Added to all three
languages for shape parity even though no current Go consumer has a
missing-timestamp case (`lambada`'s timestamps always come from a real
file's mtime) -- `*time.Time` is what makes "no value" expressible in Go at
all, the equivalent of Ruby's/Swift's native optionals.

Written by inspection per the existing sandbox limitation above (no Go
toolchain here) -- confirmed via a green CI run on GitHub Actions rather
than a locally-pasted `go vet`/`ginkgo-fd -r`, unlike every prior change
in this file. Tagged, pushed, and released:
https://github.com/woodie/humane/releases/tag/v0.9.0.

All three real consumers have since adopted it: `lambada` `2.7.0`
(deployed to the Pi, confirmed live), `scandalous` `2.7.0`, and `zouk`
`v1.11.0` (signed, notarized, and the `homebrew-zouk` cask auto-bumped via
`repository_dispatch`) -- see each repo's own `docs/COWORK.md`.

## Next up

1. `HumanSize`'s 3-significant-figure rounding rule reproduces every known
   fixture (including both cases previously cross-checked against real
   hardware) but is still an inference from a small fixture set -- worth a
   real `ByteCountFormatter` comparison across a wider range of magnitudes
   before treating it as confirmed the way the two original fixtures were.
2. `HumanSize` has no alternate unit/style options, and `TimeAgo` has no
   `.named` style (`"yesterday"`, calendar-boundary-aware) -- both left out
   deliberately, not gaps to fill without a real need.
3. `humane-ruby` issue #1 quotes ActionView's full `distance_of_time_in_words`
   table; ported only through the "1 day" row. The 2..29-day and month/year
   buckets past that are out of scope by design (see README "Scope") -- not
   a gap to fill without a real downstream need.

All three real consumers (`lambada`, `scandalous`, `zouk`) have already
adopted `v0.9.0` -- see "All three real consumers have since adopted it"
above. Nothing left open on that front.
