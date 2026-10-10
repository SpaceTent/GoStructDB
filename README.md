# GoStructDB (gsdb)

[![Go Reference](https://pkg.go.dev/badge/github.com/SpaceTent/GoStructDB/gsdb.svg)](https://pkg.go.dev/github.com/SpaceTent/GoStructDB/gsdb)
[![Go](https://github.com/SpaceTent/GoStructDB/actions/workflows/go.yml/badge.svg)](https://github.com/SpaceTent/GoStructDB/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/SpaceTent/GoStructDB)](https://goreportcard.com/report/github.com/SpaceTent/GoStructDB)
[![Release](https://img.shields.io/github/v/release/SpaceTent/GoStructDB)](https://github.com/SpaceTent/GoStructDB/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**An opinionated database library for Go. It's got some ORM, but it doesn't hide the SQL.**

GSDB maps plain Go structs to database rows using `db:"..."` struct tags. It writes the
boilerplate `INSERT` / `UPDATE` statements for you, loads query results straight into
structs, and gives you a `.Save()` that inserts or updates depending on the primary key.
You still write your own `SELECT`s.

Works with **MySQL / MariaDB** and **SQLite**.

---

## Contents

- [Why GSDB?](#why-gsdb)
- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Quick start](#quick-start)
- [Connecting](#connecting)
- [Struct tags](#struct-tags)
- [Usage](#usage)
- [Configuration](#configuration)
- [Trade-offs and caveats](#trade-offs-and-caveats)
- [Project layout](#project-layout)
- [Running the examples](#running-the-examples)
- [Status and roadmap](#status-and-roadmap)
- [Contributing](#contributing)
- [License](#license)

## Why GSDB?

While GORM is a popular choice, I believe strongly in developers understanding and writing
SQL queries directly. SQL proficiency should be a fundamental skill for all developers.
However, Go's verbose syntax can make database interactions lengthy and repetitive.

GSDB gives you some ORM features, but you still write SQL. It provides direct access to run
queries and read the results through simple `Record` and `Field` types, and it translates
database rows into structs (and structs into `INSERT`/`UPDATE` SQL) when the struct carries
`db` tags describing columns, primary keys and read defaults.

Another thing I didn't like about SQL handling in Go was making structs contain SQL types
such as `sql.NullString`. It ties your struct to a database and increases the scaffolding
you need to work with the data. With GSDB your fields stay primitive: `string`, `int`,
`bool`, `float64`, `time.Time`.

GSDB is designed to remove boilerplate and make CRUD applications easy to write, but that
ease of use comes with trade-offs (see [Trade-offs and caveats](#trade-offs-and-caveats)).
If you need fast, high-performance, memory-efficient database access, GSDB probably isn't
for you.

## Features

- **Struct ⇄ row mapping** with `db:"column=... primarykey=yes table=..."` tags.
- **Generated SQL**: `DB.Insert(v)` and `DB.Update(v)` build the SQL for a struct; you decide when to run it.
- **`DB.Save(v, pk)`**: `INSERT` when the primary key is zero, `UPDATE` otherwise.
- **Typed reads** with generics: `QueryStruct[T]` returns `[]T`, `QuerySingleStruct[T]` returns `T`.
- **Untyped reads**: `DB.Query` returns `[]Record`, and each `Field` has converters (`AsString`, `AsInt`, `AsDate`, ...).
- **No `sql.Null*` in your structs**: `NULL` strings and numbers read as zero values, and `NULL` dates are controlled per field with `readdefault`.
- **Hex-encoded string literals** in generated SQL (`X'...'`), so quotes and odd characters are written safely.
- **Query counters** to count how many statements ran during a block of code (experimental).
- **Optional SQL logging** via `log/slog`.

## Requirements

- **Go 1.26+** (the version declared in `go.mod`).
- **cgo and a C compiler** (`gcc` or `clang`) for SQLite. GSDB uses
  [`mattn/go-sqlite3`](https://github.com/mattn/go-sqlite3); with `CGO_ENABLED=0` the build
  still succeeds, but SQLite fails at runtime with *"go-sqlite3 requires cgo to work"*.
  MySQL works either way.
- **A database**: MySQL/MariaDB (via [`go-sql-driver/mysql`](https://github.com/go-sql-driver/mysql))
  or a SQLite file.

## Installation

```bash
go get github.com/SpaceTent/GoStructDB@latest
```

Then import the `gsdb` package:

```go
import "github.com/SpaceTent/GoStructDB/gsdb"
```

## Quick start

This complete program uses SQLite, so it runs with no database server. Save it as
`main.go` in a new module (`go mod init example && go get github.com/SpaceTent/GoStructDB@latest`)
and run `go run .`.

```go
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/SpaceTent/GoStructDB/gsdb"
)

// Person maps to the "people" table. The table name goes on the primary key field.
type Person struct {
	ID      int       `db:"column=id primarykey=yes table=people"`
	Name    string    `db:"column=name"`
	Email   string    `db:"column=email"`
	Active  bool      `db:"column=active"`
	Created time.Time `db:"column=created readdefault=zero"`
	Notes   string    `db:"column=notes omit=yes"` // never written to the DB
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// SQLite for the demo. For MySQL use gsdb.New("user:pass@tcp(localhost:3306)/mydb", logger, ctx).
	gsdb.NewSQLite3("quickstart.db", logger, context.Background())

	_, _, err := gsdb.DB.Execute(`CREATE TABLE IF NOT EXISTS people (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT, email TEXT, active BOOLEAN, created DATETIME)`)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	// Insert: Save does an INSERT because the primary key is zero.
	p := Person{Name: "Ada Lovelace", Email: "ada@example.com", Active: true, Created: time.Now()}
	id, _, err := gsdb.DB.Save(p, p.ID)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	fmt.Println("inserted id:", id)

	// Read one row back into the struct.
	ada, err := gsdb.QuerySingleStruct[Person]("SELECT * FROM people WHERE id = ?", id)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	fmt.Printf("loaded: %+v\n", ada)

	// Update: Save does an UPDATE because the primary key is non-zero.
	ada.Email = "ada@analytical.engine"
	if _, rows, err := gsdb.DB.Save(ada, ada.ID); err == nil {
		fmt.Println("updated rows:", rows)
	}

	// See the SQL GSDB generates without running it.
	sql, _ := gsdb.DB.Update(ada)
	fmt.Println("generated SQL:", sql)

	// Read many rows into a slice.
	people, _ := gsdb.QueryStruct[Person]("SELECT * FROM people WHERE active = ?", true)
	fmt.Println("active people:", len(people))

	// Or skip structs and work with Records and Fields.
	records, _ := gsdb.DB.Query("SELECT id, name FROM people")
	for _, r := range records {
		fmt.Println("record:", r["id"].AsInt(), r["name"].AsString())
	}
}
```

Real output from a first run (against `v0.0.14`):

```text
inserted id: 1
loaded: {ID:1 Name:Ada Lovelace Email:ada@example.com Active:true Created:2026-10-10 17:05:46 +0000 UTC Notes:}
updated rows: 1
generated SQL: UPDATE people SET name=X'416461204c6f76656c616365',email=X'61646140616e616c79746963616c2e656e67696e65',active=true,created='2026-10-10 17:05:46' WHERE id=1;
active people: 1
record: 1 Ada Lovelace
```

Note the `Created` value comes back labelled UTC: see [time zones](#trade-offs-and-caveats).
A longer walkthrough, including MySQL, is in [docs/getting-started.md](docs/getting-started.md).

## Connecting

GSDB keeps **one global connection** in `gsdb.DB`. Call one of the constructors once at
start-up:

```go
// MySQL / MariaDB. The DSN uses go-sql-driver/mysql syntax.
// The connection is opened lazily on the first query, with up to 3 ping attempts.
gsdb.New("user:password@tcp(127.0.0.1:3306)/mydb", logger, ctx)

// SQLite. The handle is set up immediately (no ping); the file is created if missing.
gsdb.NewSQLite3("app.db", logger, ctx)
```

Both take a `*slog.Logger` (used for errors and, if enabled, SQL logging) and a
`context.Context`.

## Struct tags

Every exported field that GSDB writes needs a `db` tag. Options are space-separated
`key=value` pairs:

| Option | Example | Meaning |
| --- | --- | --- |
| `column` | `column=email` | **Required.** The database column name. `Insert`/`Update` return an error if an exported field has no `column`. |
| `table` | `table=people` | The table name used by `Insert`/`Update`/`Save`. Put it on any one field; by convention the primary key. |
| `primarykey` | `primarykey=yes` | Marks the primary key. It is left out of `INSERT` and used as the `WHERE` clause for `UPDATE`. |
| `omit` | `omit=yes` | Never written by `Insert`/`Update`. Still filled by `QueryStruct` if the query returns that column. |
| `readdefault` | `readdefault=zero` | `time.Time` only. What to return when the column is `NULL`: `now` (current time), `zero` (`0001-01-01 00:00:00`), or `null` (zero on read, **and** a zero time is written back as SQL `NULL`). |

Unexported fields are ignored.

**`readdefault` and `NULL` dates.** If no `readdefault` is set, a `NULL` date is read as
`time.Now()`. Early versions of GSDB did this and some applications rely on it, so it
remains the default. For new code use `readdefault=zero` (or `readdefault=null` if you also
want zero times saved as `NULL`).

**Supported field types**

| Go type | Insert / Update | QueryStruct |
| --- | --- | --- |
| `int`, `int8`–`int64`, `uint`, `uint8`–`uint64` | ✅ | ✅ |
| `float32`, `float64` | ✅ | ✅ |
| `bool` | ✅ (written as `true`/`false`) | ✅ |
| `string` | ✅ (written as hex `X'...'`) | ✅ |
| `time.Time` | ✅ (written as `'YYYY-MM-DD HH:MM:SS'`) | ✅ |
| `[]byte` | ❌ | ✅ |

Pointers and `sql.Null*` types are not supported.

## Usage

### Executing statements

```go
lastInsertedID, rowsAffected, err := gsdb.DB.Execute("DELETE FROM people WHERE id = ?", 42)
```

`Execute` passes parameters to `database/sql`, so use `?` placeholders for values.

### Generating INSERT and UPDATE SQL

`Insert` and `Update` only **build** SQL from a struct; run it with `Execute`.

```go
sqlQuery, err := gsdb.DB.Insert(person) // INSERT INTO people(name,email,...) VALUES (...);
// or: sqlQuery, err := gsdb.DB.Update(person) // UPDATE people SET name=...,email=... WHERE id=1;
if err != nil {
	return err
}
lastInsertedID, rowsAffected, err := gsdb.DB.Execute(sqlQuery)
```

Strings are written as hex literals, which keeps quotes and unusual characters intact:

```sql
INSERT INTO Test(name,dtadded,status) VALUES (X'54657374','2025-12-25 15:29:25',1);
```

### Save

`Save` takes the struct and its primary key value. If the key is the zero value it runs an
`INSERT`; otherwise an `UPDATE`.

```go
id, rows, err := gsdb.DB.Save(person, person.ID)
```

### Reading into structs

`QueryStruct` and `QuerySingleStruct` are package functions (Go does not allow generic
methods). Columns are matched to fields by the `column` tag.

```go
people, err := gsdb.QueryStruct[Person]("SELECT * FROM people WHERE active = ?", true)

person, err := gsdb.QuerySingleStruct[Person]("SELECT * FROM people WHERE id = ?", id)
// If no row matches, person is the zero value and err is nil.
```

### Records and Fields

`DB.Query` returns `[]Record`, where a `Record` is `map[string]Field`. A `Field` wraps the
raw value and offers converters:

```go
records, err := gsdb.DB.Query("SELECT id, name, created FROM people")
for _, r := range records {
	id := r["id"].AsInt()                  // also AsInt64, AsUInt64, AsFloat, AsBool
	name := r["name"].AsString()           // also AsByte
	created := r["created"].AsDate("zero") // "now", "zero", or a "2006-01-02 15:04:05" fallback
	fmt.Println(id, name, created)
}
```

`RecordInsert(record, table)` and `RecordUpdate(record, table, column, value)` write a
`Record` back without needing a struct.

### Counters (experimental)

Count how many statements run between two points in your code:

```go
gsdb.DB.StartCounter("request")
// ... lots of DB calls ...
count := gsdb.DB.GetCounter("request")
```

Every `Execute` and `Query` increments all active counters.

## Configuration

These package-level variables can be set before use:

| Variable | Default | Effect |
| --- | --- | --- |
| `gsdb.ShowSQL` | `false` | Log every statement run through `Execute` (with last insert ID and rows affected) at `INFO`. |
| `gsdb.ColumnWarnings` | `false` | Log a warning when a query returns a column that has no matching struct field. |
| `gsdb.MaxOpenConns` | `50` | Max open connections (MySQL connections opened by `gsdb.New`). |
| `gsdb.MaxIdleConns` | `50` | Max idle connections (MySQL). Idle connections are closed after 5 minutes. |

```go
gsdb.ShowSQL = true
gsdb.ColumnWarnings = true
```

## Trade-offs and caveats

- **Reflection and memory.** GSDB uses reflection for type conversion and loads the whole
  result set into memory. Paginate large queries.
- **One global connection.** `gsdb.DB` is package-level; you can't hold two databases at once.
- **Time zones.** Times are written as `YYYY-MM-DD HH:MM:SS` with no zone, and read back
  with no zone (Go treats that as UTC). Store and compare times in UTC to avoid surprises.
- **Generated SQL embeds values.** `Insert`/`Update`/`Save` put values directly into the SQL
  text (strings hex-encoded, numbers and dates formatted). Use `Execute`/`Query` with `?`
  parameters for anything you write by hand.
- **Databases.** Only the MySQL and SQLite drivers are registered. The generated SQL (hex
  literals, `true`/`false`) targets those two.
- **`InsertMany` is not ready.** It currently joins rows without commas, which produces
  invalid SQL for more than one row. Use `Insert` in a loop for now.

## Project layout

```text
.
├── gsdb/               # The library
│   ├── New.go          # Database type, New / NewSQLite3, connection pool settings
│   ├── Execute.go      # DB.Execute
│   ├── Query.go        # DB.Query -> []Record
│   ├── QueryStruct.go  # QueryStruct[T] / QuerySingleStruct[T]
│   ├── Insert.go       # DB.Insert, InsertMany
│   ├── Update.go       # DB.Update
│   ├── Save.go         # DB.Save (insert or update)
│   ├── Record.go       # Record type, RecordInsert / RecordUpdate
│   ├── Fields.go       # Field type and As* converters
│   ├── counters.go     # Statement counters
│   ├── common.go       # Tag parsing and SQL builders
│   └── *_test.go       # Unit tests (sqlmock) and SQLite integration tests
├── examples/           # Runnable examples, called from main.go
├── main.go             # Runs all the examples
└── docs/               # Longer guides
```

How the pieces fit together is described in [docs/architecture.md](docs/architecture.md).

## Running the examples

The `examples/` directory has runnable programs for inserts, counters, missing columns,
date defaults, `Save`, SQLite and concurrent inserts. See
[examples/README.md](examples/README.md) for how to run them.

## Status and roadmap

GSDB is pre-1.0 (`v0.0.x`), so the API may still change between releases. See [CHANGELOG.md](CHANGELOG.md) and the
[releases page](https://github.com/SpaceTent/GoStructDB/releases).

Items noted in the code as not finished yet:

- Counters are experimental.
- Wider type coverage in `Query` (marked `TODO` in the source).
- `InsertMany` (see caveats above).

## Contributing

Contributions are welcome. Please read [CONTRIBUTING.md](CONTRIBUTING.md) and our
[Code of Conduct](CODE_OF_CONDUCT.md). To report a security issue, see
[SECURITY.md](SECURITY.md).

## License

GoStructDB is released under the [MIT License](LICENSE).
