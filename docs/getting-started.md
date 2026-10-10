# Getting started with GSDB

This guide takes you from an empty folder to a working program that inserts, reads,
updates and deletes rows with GSDB. It uses SQLite first (no server needed), then shows the
changes for MySQL.

## 1. Prerequisites

- Go 1.26 or newer: `go version`
- A C compiler for SQLite support (cgo): `gcc --version` or `clang --version`
  - macOS: `xcode-select --install`
  - Debian/Ubuntu: `sudo apt install build-essential`

## 2. Create a module

```bash
mkdir gsdb-demo && cd gsdb-demo
go mod init example.com/gsdb-demo
go get github.com/SpaceTent/GoStructDB@latest
```

## 3. Describe your table as a struct

```go
type Person struct {
	ID      int       `db:"column=id primarykey=yes table=people"`
	Name    string    `db:"column=name"`
	Email   string    `db:"column=email"`
	Active  bool      `db:"column=active"`
	Created time.Time `db:"column=created readdefault=zero"`
	Notes   string    `db:"column=notes omit=yes"`
}
```

- `column=` names the column. Every exported field needs one.
- `table=` (on one field, usually the primary key) names the table for generated SQL.
- `primarykey=yes` is skipped on `INSERT` and becomes the `WHERE` clause on `UPDATE`.
- `omit=yes` is never written.
- `readdefault=zero` returns the zero time when the column is `NULL` (without it, GSDB
  returns `time.Now()` for backward compatibility).

The full tag reference is in the [README](../README.md#struct-tags).

## 4. Write the program

Copy the [Quick start program](../README.md#quick-start) into `main.go`. It:

1. Connects with `gsdb.NewSQLite3("quickstart.db", logger, ctx)`.
2. Creates the `people` table with `gsdb.DB.Execute`.
3. Inserts a `Person` with `gsdb.DB.Save(p, p.ID)` (ID is 0, so it's an `INSERT`).
4. Loads it with `gsdb.QuerySingleStruct[Person]`.
5. Changes the email and calls `Save` again (ID is now set, so it's an `UPDATE`).
6. Prints the SQL `gsdb.DB.Update` would generate.
7. Lists rows with `gsdb.QueryStruct[Person]` and with the untyped `gsdb.DB.Query`.

## 5. Run it

```bash
go run .
```

```text
inserted id: 1
loaded: {ID:1 Name:Ada Lovelace Email:ada@example.com Active:true Created:2026-10-10 17:05:46 +0000 UTC Notes:}
updated rows: 1
generated SQL: UPDATE people SET name=X'416461204c6f76656c616365',email=X'61646140616e616c79746963616c2e656e67696e65',active=true,created='2026-10-10 17:05:46' WHERE id=1;
active people: 1
record: 1 Ada Lovelace
```

Run it again and you'll get `inserted id: 2`, because the database file persists.

## 6. Turn on SQL logging while you develop

```go
gsdb.ShowSQL = true        // log every Execute
gsdb.ColumnWarnings = true // warn when a column has no matching struct field
```

With `ShowSQL` on, each statement is logged through your `slog.Logger`:

```text
level=INFO msg="INSERT INTO people(name,created) VALUES (X'78','2026-10-10 17:22:20');" lastid=1 "rows effected"=1
```

## 7. Delete rows

There's no `Delete` helper; write the SQL and use parameters:

```go
_, rows, err := gsdb.DB.Execute("DELETE FROM people WHERE id = ?", ada.ID)
```

## 8. Switching to MySQL / MariaDB

Only the connection line and the DDL change.

```go
gsdb.New("user:password@tcp(127.0.0.1:3306)/gsdb_demo", logger, context.Background())
```

```sql
CREATE TABLE people (
  id      INT AUTO_INCREMENT PRIMARY KEY,
  name    VARCHAR(255),
  email   VARCHAR(255),
  active  BOOLEAN,
  created DATETIME NULL
);
```

Notes:

- The DSN uses [go-sql-driver/mysql syntax](https://github.com/go-sql-driver/mysql#dsn-data-source-name).
- `gsdb.New` doesn't connect straight away. The first `Execute`/`Query` opens the pool and
  pings the server (up to 3 attempts, 500 ms apart).
- Pool size comes from `gsdb.MaxOpenConns` / `gsdb.MaxIdleConns` (default 50). Set them
  before the first query.

A throwaway MySQL server for local testing, if you use Docker:

```bash
docker run --rm -d --name gsdb-mysql -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=secret -e MYSQL_DATABASE=gsdb_demo mysql:8
# DSN: root:secret@tcp(127.0.0.1:3306)/gsdb_demo
```

## Next steps

- Try the runnable programs in [examples/](../examples/README.md).
- Read [how GSDB works](architecture.md).
- Check the [trade-offs and caveats](../README.md#trade-offs-and-caveats) before using it on large tables.
