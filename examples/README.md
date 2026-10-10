# Examples

Each file in this directory is a small function showing one GSDB feature. `main.go` in the
repo root calls them all, in this order:

| Function | File | Shows |
| --- | --- | --- |
| `BasicInsert` | `basic-Insert.go` | Build an `INSERT` from a struct with `DB.Insert`, run it with `DB.Execute` |
| `BasicInsertWithCounter` | `basic-insert-withcounter.go` | `StartCounter` / `GetCounter` around 200 queries |
| `ConcurrentInserts` | `Concurrent-Inserts.go` | 1,000 goroutines inserting, reading back and deleting |
| `InsertDateDefault` | `Insert-DateDefault.go` | `readdefault=now`, `zero`, `null` and the default for `NULL` dates |
| `MissingColumns` | `Missing-Columns.go` | Struct/table mismatches and `gsdb.ColumnWarnings` |
| `Lite3` | `SQLite3.go` | `NewSQLite3` against `test.db`, creating the `Test` table |
| `SaveReadDefaultNull` | `Save-DateDefault.go` | `DB.Save` insert then update with `readdefault=null` |

All examples connect with `NewSQLite3` using the `DSN` environment variable (loaded from a
`.env` file if present). The MySQL call `New(DSN, ...)` is left commented out in each file.

## Run them against SQLite

From the repository root:

```bash
# 1. Create the Test table the examples expect
#    (Lite3 creates it too, but it runs after the others).
sqlite3 test.db "CREATE TABLE IF NOT EXISTS Test (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, dtadded DATETIME DEFAULT CURRENT_TIMESTAMP, status INTEGER);"

# 2. Run every example
DSN=test.db go run main.go
```

No `sqlite3` CLI? Run `DSN=test.db go run main.go` twice; the first run creates the table
(via `Lite3`) and the second run works.

You can also put the DSN in a file instead of the environment. Copy `examples/env-dist`:
`BasicInsert` reads `examples/.env`, and the other examples read `.env` in the current
directory, so `DSN=...` in the environment is the simplest option.

### What to expect

You'll see `slog` lines such as:

```text
level=INFO msg="Item with ID 1 was inserted. 1 rows were affected"
level=INFO msg="Number of SQL Queries" Counter=201
level=INFO msg=Date "Date Should be NULL/zero"=0001-01-01T00:00:00.000Z
level=INFO msg="Save insert rows: 1, save update rows: 1"
```

`ConcurrentInserts` is written with MySQL in mind. On SQLite, 1,000 concurrent writers
produce many `database is locked` errors (and some `ID mismatch` lines); that is SQLite
write locking, not a GSDB failure.

`MissingColumns` deliberately prints a `WARN Unknown type col=status ...` line to show
`ColumnWarnings`.

## Run them against MySQL

1. In each example, swap `MySQL.NewSQLite3(DSN, ...)` for the commented-out
   `MySQL.New(DSN, ...)`.
2. Create the table:

   ```sql
   CREATE TABLE Test (
     id INT AUTO_INCREMENT PRIMARY KEY,
     name VARCHAR(255),
     dtadded DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
     status INT
   );
   ```

3. Run with a go-sql-driver DSN:

   ```bash
   DSN='user:password@tcp(127.0.0.1:3306)/gsdb_examples' go run main.go
   ```

Some examples query `test` (lower case) while the table is `Test`. That is fine on SQLite
and on MySQL with case-insensitive table names (the default on macOS and Windows). On Linux
MySQL servers with `lower_case_table_names=0`, create the table as `test` too or adjust the
queries.
