# Go StructDB (GSDB)

GSDB is a lightweight database library for Go.  
It provides basic ORM-style features but keeps SQL explicit.

Unlike full ORMs such as GORM, GSDB requires you to write SQL queries directly while reducing repetitive boilerplate code.

---

## Philosophy

- SQL should be written and understood by developers.
- Structs should use Go primitives (`string`, `time.Time`) instead of database-specific types like `sql.Null*`.
- Convenience methods should simplify CRUD operations without hiding the underlying SQL.

---

## Installation

```bash
go get github.com/SpaceTent/GoStructDB/gsdb
```

## Connecting
```go
db := gsdb.New(dataSource, logger, ctx)
```

## Executing Queries

```go
lastInsertedID, rowsAffected, err := db.Execute(sqlQuery, params...)

```

## Insert and Update from Structs
GSDB can build SQL statements from structs with db.Insert and db.Update.

```go
type InsertPerson struct {
    Id      int       `db:"column=id primarykey=yes table=Test"`
    Name    string    `db:"column=name"`
    Dtadded time.Time `db:"column=dtadded"`
    Status  int       `db:"column=status"`
    Ignored int       `db:"column=ignored omit=yes"`
}

entry := InsertPerson{
    Name:    "Test",
    Dtadded: time.Now(),
    Status:  1,
}

sqlQuery, err := db.Insert(entry)
if err != nil {
    log.Error(err.Error())
    return
}

lastInsertedID, rowsAffected, err := db.Execute(sqlQuery)
```
### Resulting SQL:
```sql
INSERT INTO Test(name,dtadded,status)
VALUES (X'54657374','2025-12-25 15:29:25',1);
```

## Query Structs

### QueryStruct
Returns a slice of structs
```go
people, err := db.QueryStruct[InsertPerson]("SELECT * FROM Test")
```

### QuerySingleStruct
Returns a single struct

```go
person, err := db.QuerySingleStruct[InsertPerson]("SELECT * FROM Test WHERE id = ?", id)
```

## Handling NULL Dates
- readdefault=now → NULL becomes time.Now(). (Legacy, not recommended.)
- readdefault=zero → NULL becomes Go’s zero time (0001-01-01 00:00:00). (Recommended.)

## Save
`db.Save` decides whether to insert or update based on the primary key
```go
entry := InsertPerson{
    Name:    "Test",
    Dtadded: time.Now(),
    Status:  1,
}

db.Save(entry, 0)
```

if `Id == 0`, it runs INSERT, otherwise, it runs an UPDATE

## Counters (experimental)
Track how many queries are executed between checkpoints
```go
db.StartCounter("test")
// multiple db calls
count := db.GetCounter("test")
```

## Column Warnings
If struct fields don't match the query result, GSDB raises a warning
```go
db.ColumnWarnings = true
p, _ := db.QuerySingleStruct[InsertPerson]("SELECT * FROM Test WHERE id = ?", id)
```

## Trade-offs
- Uses reflection for type conversions
- Loads complete record sets into memory
- Not intended for high-performance or large-scale data processing

## Summary
GSDB reduces Go’s database boilerplate while keeping SQL explicit.
It is best suited for small to medium CRUD applications where clarity and simplicity matter more than raw performance.