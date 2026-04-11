.PHONY: setup test lint

# first thing after cloning
setup:
	git config core.hooksPath .githooks
	@echo "hooks activated"

test:
	go test -race ./...

lint:
	golangci-lint run ./...

