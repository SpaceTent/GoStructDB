# Go StructDB - gsdb

Opinionated Database Extraction Library for Go, It's got some ORM, but does not hide the SQL. 

While GORM is a popular choice, I believe strongly in developers understanding and writing SQL queries directly. 

SQL proficiency should be a fundamental skill for all developers. However, Go's verbose syntax can make database
interactions quite lengthy and repetitive.

GSDB give you some ORM features but you still need to write SQL. 

GSDB works for me, as it provides direct access to run queries and access the results with simple "record" and "field" types.
It also provides Database to Struct transformations. If the struct has the extra "db" tags to define columns, primary keys read default conditions, GSDB will translate DB resources into Structs and provide methods such as .Save(); Making it easy to define a struct, load it with data and just call a DB.Save to save it to the database

GSDB is designed to abstract a lot of the boilerplate DB code away from you. It makes writing CRUD applications easier, but, with that ease of use, comes trade-offs.
GSDB uses reflection to determine many of the type conversions. It also holds complete record sets in memory. 
If you need fast, high performance and efficient DB access, GSDB probably isn't for you.

Another thing I didn't like about SQL handling in Go, was making structs contain SQL data types such as sql.Null.
To me this seems lazy. It ties your struct to a database, increasing the boilerplate and scaffolding you need to work with the data. To me, the data is either strings or dates. Then the struct is either string or time.Time.

I'm keen to keep primitives, as primitive as possible. 

### Connecting

To start a new DB connection, include GSDB library and declare a new connection

```go
gsdb.New(DataSource, StructuredLoggingHandler, context)
```

### Executing Queries

```go
lastInsertedID, rowsAffected, err := gsdb.DB.Execute(sqlQuery,parameters...)
```
### Making Insert and Update SQL

gsdb provides 2 functions to make SQL statements from structs, DB.Insert and DB.Update. If you pass in any variable from a struct, the library will build the SQL necessary for an INSERT or UPDATE statement.

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
	sqlQuery, err := MySQL.DB.Insert(entry)
	if err != nil {
		l.Error(err.Error())
		return
	}
	
	lastInsertedID, rowsAffected, err := MySQL.DB.Execute(sqlQuery)
	if err != nil {
		l.Error(err.Error())
	}
	l.Info(fmt.Sprintf("Item with ID %d was inserted. %d rows were affected", lastInsertedID, rowsAffected))
```

### Hex conversion of strings

All strings are converted to HEX. This ensures even the most challenging characters are written to the DB as it also makes SQL injections attacks difficult.  

```sql
INSERT INTO Test(name,dtadded,status) VALUES (X'54657374','2025-12-25 15:29:25',1);
```

### QueryStruct

QueryStruct will return a slice of type T containing all the data.

### QuerySingleStruct

This works the same as QueryStruct, but returns T, not a slice of type T. 
When returning datetimes from the database, there are some extra options to determine how you can view NULL or EMPTY dates in the database. "readdefault=now" or "readdefault=zero"

In early versions of GSDB, if a date field in the database was NULL, the library returned time.Now() instead. That behavior wasn’t really correct, but some applications had already built business logic around the idea that a NULL date meant “current time.” To preserve compatibility, GSDB kept this as the default.

However, the recommended approach is to use the tag readdefault=zero, which tells GSDB to return the Go “zero time” (0001-01-01 00:00:00) for NULL dates. This is usually the proper way to represent missing or unset timestamps.

### Save

If the primary key is zero, then an Insert is executed, otherwise it's an Update. 

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
    
    gsdb.New(DataSource, StructuredLoggingHandler, context)
    
    gsdb.Save(entry,0)
    
 ```   

### Counters

You can start a counter anywhere in your call code, and then call the getCounter functions to see how many SQL statements have happened since that counter was started. 

This is experimental at this stage. 

```go
	MySQL.DB.StartCounter("test")
	
	// Loads of DB calls 
	
	count := MySQL.DB.GetCounter("test)
```

### ColumnWarnings

If the struct doesn't match the SQL returned or generated.  The library will spit out a warning that there is a mismatch, and then ignore it.  You can turn of this warning with 

```go
    MySQL.ColumnWarnings = true
    P, _ := MySQL.QuerySingleStruct[InsertPerson]("select * from Test WHERE ID = ?", lastInsertedID)
```
If there are any fields defined in InsertPerson that don't match the record set, a warning is raised. 


### ShowSQL

