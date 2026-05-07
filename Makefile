.PHONY: fmt fmt-check

# Format all Go files in the repo (same style as gofmt pre-commit).
fmt:
	gofmt -w .

# Fail if any .go file is not gofmt-clean (for CI or manual check).
fmt-check:
	@test -z "$$(gofmt -l .)" || (echo "gofmt needed on:" && gofmt -l . && exit 1)
