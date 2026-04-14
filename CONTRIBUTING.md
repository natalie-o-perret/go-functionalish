# Contributing

## Setup

```sh
git clone https://github.com/natalie-o-perret/gof.git
cd gof
make setup   # installs the commit-msg hook
```

Requires Go 1.24+ and [golangci-lint](https://golangci-lint.run/welcome/install/).

## Workflow

```sh
make test      # go test -race ./...
make lint      # golangci-lint run ./...
make lint-md   # markdownlint on all .md files
```

All three run in CI on every PR. A PR must be green before merging.

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

The commit-msg hook (installed by `make setup`) enforces this locally. The PR title is also linted in CI.

`BREAKING CHANGE:` in the commit footer signals a semver major bump to release-please.

## Pull requests

- One logical change per PR.
- Add or update tests for any behaviour change.
- No external dependencies - the library intentionally has none.

## Design constraints

- **No reflection, no `interface{}`** - pure generics only.
- **Type-transforming functions are package-level** (Go methods cannot introduce new type parameters).
- **`seq` operations are lazy** - wrapping iterators, no materialisation until a terminal is called.
- **F# naming conventions** - `Bind` not `FlatMap`, `Map` not `Select`, `Some`/`None` not `Just`/`Nothing`.
