# Changelog

## [2.0.0](https://github.com/natalie-o-perret/go-functionalish/compare/v1.5.0...v2.0.0) (2026-04-26)


### ⚠ BREAKING CHANGES

* FlatMap renamed to Bind in option, result, and validation

### Features

* add option.Map2, result.Zip, update README ([#10](https://github.com/natalie-o-perret/go-functionalish/issues/10)) ([c5f7694](https://github.com/natalie-o-perret/go-functionalish/commit/c5f769483901e8d275c7f6d11e860abd5ef37eef))
* add pipe compose, curried helpers, ToSlice rename, benchmarks, fix ci ([#1](https://github.com/natalie-o-perret/go-functionalish/issues/1)) ([b6c3258](https://github.com/natalie-o-perret/go-functionalish/commit/b6c3258a0e288a7e554ee1200583ac57f55194a1))
* add pseq package — parallel sequence operations ([#12](https://github.com/natalie-o-perret/go-functionalish/issues/12)) ([6fd094a](https://github.com/natalie-o-perret/go-functionalish/commit/6fd094a6ada322aaf740b3306373d4587ee97d46))
* add tests, ci, govulncheck, release workflows and linter config ([8381b73](https://github.com/natalie-o-perret/go-functionalish/commit/8381b7326bb4f1ab52fe5bfc00d442f6b507eef8))
* add tuple package ([#11](https://github.com/natalie-o-perret/go-functionalish/issues/11)) ([50817bb](https://github.com/natalie-o-perret/go-functionalish/commit/50817bbf11ae076587385e17c194c1b7a25302a6))
* add Unfold, Partition, OfMap, OfOption, RangeStep, CountByKey, Interleave, Flatten, OrElse, Tee/TeeErr, Tap ([#9](https://github.com/natalie-o-perret/go-functionalish/issues/9)) ([6cb5728](https://github.com/natalie-o-perret/go-functionalish/commit/6cb57286ce086931df927293a9933ca2f8a90561))
* add validation package for applicative error accumulation ([#4](https://github.com/natalie-o-perret/go-functionalish/issues/4)) ([73c18a2](https://github.com/natalie-o-perret/go-functionalish/commit/73c18a20f91de58cfce9fe30be2d9c2b6cf12761))
* add validation package, Bind rename, partial type inference ([#5](https://github.com/natalie-o-perret/go-functionalish/issues/5)) ([23d331e](https://github.com/natalie-o-perret/go-functionalish/commit/23d331ed1cb2cedda6fa50e70ae68970eb6104f5))
* initial commit — enum, option, result, pipe packages ([0f95b47](https://github.com/natalie-o-perret/go-functionalish/commit/0f95b470a0015ceec689a991c9fe8f56e957a991))


### Documentation

* add root package doc for pkg.go.dev ([#3](https://github.com/natalie-o-perret/go-functionalish/issues/3)) ([2e85eba](https://github.com/natalie-o-perret/go-functionalish/commit/2e85eba134ae98456f48e0a87c5c113ce06305bc))
* rename PipeN/ComposeN to PipeEndoN/ComposeEndoN in README ([#6](https://github.com/natalie-o-perret/go-functionalish/issues/6)) ([1b4b14a](https://github.com/natalie-o-perret/go-functionalish/commit/1b4b14aa99334a4bfc7a3724e617c7b534655b85))
