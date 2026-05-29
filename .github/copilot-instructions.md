# Copilot code review instructions for go-functionalish

## General Go guidelines

- Flag any use of `interface{}` or `any` where a concrete typed alternative exists.
- Prefer explicit error handling; every `error` return must be checked.
- Ensure all exported functions, types, and constants have doc comments following the `// Name ...` convention.
- Use table-driven tests; flag test functions that lack cases for edge values (nil, empty, single-element).

## Functional / generics API design

- Warn if a newly added function could be expressed as a composition of existing primitives in the package.
- Ensure generic type constraints are as narrow as needed — prefer `comparable` or interface constraints over unconstrained `any`.
- Flag mutable closures inside pure functions (captures that write to outer state).
- Prefer returning a new value over mutating a receiver; flag in-place mutation unless clearly documented.

## Test coverage

- Every exported function must have at least one test exercising the happy path and one exercising a boundary/error case.
- Benchmarks should be present for hot-path functions operating on sequences.
- Flag calls to `t.Skip` without a corresponding tracking issue reference.

## Dependency hygiene

- This library targets zero non-test dependencies beyond the standard library and `samber/lo`; flag any new runtime dependency.
- Warn on `go.sum` changes without a corresponding `go.mod` change.
