# Contributing to GoStructDB

Thanks for your interest in improving GSDB! Bug reports, docs fixes, tests and features are
all welcome.

By taking part you agree to follow our [Code of Conduct](CODE_OF_CONDUCT.md).

## Before you start

- **Bugs:** search [existing issues](https://github.com/SpaceTent/GoStructDB/issues) first,
  then open a bug report with a minimal reproduction (struct, table schema, GSDB calls).
- **Features:** open a feature request to discuss the idea before writing a lot of code.
  GSDB is deliberately opinionated (SQL stays visible, structs stay primitive), so it's
  worth checking the idea fits.
- **Security issues:** don't open a public issue; follow [SECURITY.md](SECURITY.md).

## Development setup

You need:

- Go **1.26+** (see `go.mod`)
- A C compiler (`gcc`/`clang`) with cgo enabled, for `mattn/go-sqlite3`
- Optional: `goimports` and `govulncheck` (used by `make fmt` / `make commit`)

```bash
git clone https://github.com/SpaceTent/GoStructDB.git
cd GoStructDB
go build ./...
go test ./...
```

You don't need a database server to run the tests: unit tests use
[`go-sqlmock`](https://github.com/DATA-DOG/go-sqlmock) and the integration tests use a
temporary SQLite file.

Optional tools:

```bash
go install golang.org/x/tools/cmd/goimports@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
```

## Making changes

1. Fork the repo and create a branch from `main`:
   `git checkout -b fix/short-description` (or `feat/...`, `docs/...`, `test/...`).
2. Make your change. Keep it focused: one fix or feature per pull request.
3. Add or update tests in `gsdb/*_test.go`. Table-driven tests with `testify/assert`
   match the existing style.
4. Format, vet and test:

   ```bash
   gofmt -s -w .
   goimports -w .      # if installed
   go vet ./...
   go test ./...       # or: make test (writes profile.cov with coverage)
   ```

5. Update `README.md` / `docs/` if behaviour or the API changed, and add a line under
   **Unreleased** in [CHANGELOG.md](CHANGELOG.md).

### Make targets

| Target | What it does |
| --- | --- |
| `make test` | `go test -v -coverprofile=profile.cov ./...` |
| `make fmt` | `go mod tidy`, `goimports`, `gofmt -s` |
| `make run` | Formats, then runs `main.go` (all examples; see [examples/README.md](examples/README.md)) |
| `make commit` / `make release` | Maintainer shortcuts that commit, tag and **push**. Contributors should use normal git and pull requests instead. |

## Commit messages

Short, imperative subject lines. A [Conventional Commits](https://www.conventionalcommits.org/)
style prefix is preferred and already used in the history:

```text
fix: handle NULL dates in QuerySingleStruct
feat: support uint primary keys in Save
docs: document readdefault=null
test: add Update tests for bool fields
chore: bump go-sql-driver/mysql
```

## Pull requests

- Open the PR against `main` and fill in the template.
- CI (`.github/workflows/go.yml`) runs `go test -v ./...` on every pull request; it must pass.
- Keep PRs small enough to review. Explain *why* as well as *what*.
- Maintainers may squash commits when merging.

## Releases (maintainers)

Releases are tagged `v0.0.N`. Pushing a `v*` tag runs `.github/workflows/release.yml`, which
creates a GitHub release with generated notes. Move the **Unreleased** entries in
`CHANGELOG.md` under the new version when tagging.
