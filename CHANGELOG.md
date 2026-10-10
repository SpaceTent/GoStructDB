# Changelog

All notable changes to this project are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this
project uses [Semantic Versioning](https://semver.org/) tags (pre-1.0, so minor API changes
can happen in any `v0.0.x` release). Earlier entries were reconstructed from the git history
and GitHub releases.

## [Unreleased]

### Added
- LICENSE (MIT), CONTRIBUTING, CODE_OF_CONDUCT, SECURITY, issue and PR templates.
- Expanded README, `docs/getting-started.md`, `docs/architecture.md` and `examples/README.md`.

## [v0.0.14] - 2026-09-08

### Changed
- Connection pool size is now configurable through `gsdb.MaxOpenConns` and
  `gsdb.MaxIdleConns` (default 50 each, previously fixed at 25).

## [v0.0.13] - 2026-03-06

### Added
- `Save` example for `readdefault=null` dates (`examples/Save-DateDefault.go`).

### Changed
- Examples now run against SQLite (`NewSQLite3`) by default.

## [v0.0.12] - 2026-03-04

### Added
- `readdefault=null` for `time.Time` fields: zero times are written as SQL `NULL` on
  insert/update, and `NULL` is read back as the zero time.

## [v0.0.11] - 2026-02-25

### Added
- Tests for `Execute`, `Record` insert/update and `QueryStruct`.

### Changed
- README re-worded and typos fixed.
- Dependency updates.

## [v0.0.7] - 2025-05-14

### Added
- First GitHub release: `New`, `NewSQLite3`, `Execute`, `Query`, `QueryStruct`,
  `QuerySingleStruct`, `Insert`, `Update`, `Save`, `Record`/`Field` helpers, counters,
  and a GitHub Actions release workflow.

Tags `v0.0.0`–`v0.0.6` and `v0.0.8`–`v0.0.10` exist in git but have no GitHub release
or separate notes; their changes are folded into the nearest release above.

[Unreleased]: https://github.com/SpaceTent/GoStructDB/compare/v0.0.14...HEAD
[v0.0.14]: https://github.com/SpaceTent/GoStructDB/compare/v0.0.13...v0.0.14
[v0.0.13]: https://github.com/SpaceTent/GoStructDB/compare/v0.0.12...v0.0.13
[v0.0.12]: https://github.com/SpaceTent/GoStructDB/compare/v0.0.11...v0.0.12
[v0.0.11]: https://github.com/SpaceTent/GoStructDB/compare/v0.0.7...v0.0.11
[v0.0.7]: https://github.com/SpaceTent/GoStructDB/releases/tag/v0.0.7
