.PHONY: lint test check

# lint and test are always verbose. check is terse (dots on pass, full
# detail on any failure/error) -- matching Ruby's/Swift's own lint/test/check
# split in this family.

lint:
	golangci-lint run

# Verbose on purpose -- ginkgo-fd's documentation-style output, the Go
# equivalent of Ruby's `rspec -fd` / Swift's `swift test | xctidy`.
test:
	ginkgo-fd -r

# Terser than `test` on purpose: plain ginkgo's default reporter prints a
# dot per passing spec and suppresses per-spec chatter, but always prints
# full detail for any failure.
check: lint
	ginkgo -r
