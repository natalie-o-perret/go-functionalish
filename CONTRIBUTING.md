# Contributing

## Setup

```sh
git clone https://github.com/natalie-o-perret/go-functionalish.git
cd go-functionalish
git config core.hooksPath .githooks   # installs the commit-msg hook
```

Requires Go 1.24+ and [golangci-lint](https://golangci-lint.run/welcome/install/).

## Workflow

```sh
go test -race ./...          # run tests
golangci-lint run ./...      # lint Go files
```

All checks run in CI on every PR. A PR must be green before merging.

## Commit messages

Commits must follow [Conventional Commits](https://www.conventionalcommits.org/):

```text
type(scope?): short description
```

Allowed types: `feat` `fix` `docs` `style` `refactor` `perf` `test` `build` `ci` `chore` `revert`

```text
feat: add Zip to seq package
fix(option): handle None in UnwrapOrElse
docs: update README pipeline example
refactor(result): simplify Bind internals
```

The commit-msg hook (installed above) enforces this locally. The PR title is also linted in CI.

`BREAKING CHANGE:` in the commit footer signals a semver major bump.
