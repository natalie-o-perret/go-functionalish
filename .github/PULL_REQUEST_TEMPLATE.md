## Summary

Tightens the linter setup so dead code and warnings are caught in CI before anything lands on main.

## Related issue

N/A

## Changes

- Added `deadcode` linter -- finds unreachable top-level declarations that `unused` alone misses
- Added `nilerr` -- flags the classic "return nil, nil" when the error check is wrong
- Added `wastedassign` -- catches assignments whose result is never read
- PR template checklist now has two explicit blocking items for dead code and warnings

## Checklist

- [x] No code changes, no tests needed
- [x] `golangci-lint run ./...` not applicable (config + Markdown only)

## Notes for reviewers

All three new linters are zero-tolerance: any finding is a CI error and blocks merge. If a specific false-positive comes up, add a targeted exclusion in `.golangci.yml` rather than disabling the linter wholesale.
