.PHONY: setup test lint lint-md

# first thing after cloning
setup:
	git config core.hooksPath .githooks
	@echo "hooks activated"

test:
	go test -race ./...

lint:
	golangci-lint run ./...

lint-md:
	npx --yes markdownlint-cli2 "**/*.md"

