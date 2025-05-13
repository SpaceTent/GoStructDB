package examples

import (
    "context"
    "fmt"
    "log/slog"
    "os"
    "time"
    
    gsdb "GoStructDB/gsdb"
)

func Lite3() {
    
    textHandler := slog.NewTextHandler(os.Stdout, nil)
    l := slog.New(textHandler)
    
    gsdb.NewSQLite3("test.db", l, context.Background())
    
    sql := "CREATE TABLE IF NOT EXISTS Test (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, dtadded DATETIME DEFAULT CURRENT_TIMESTAMP, status INTEGER);"
    
    gsdb.DB.Execute(sql)
    
    type InsertPerson struct {
        Id      int       `db:"column=id primarykey=yes table=Test"`
        Name    string    `db:"column=name"`
        Dtadded time.Time `db:"column=dtadded"`
        Status  int       `db:"column=status"`
        Ignored int       `db:"column=ignored omit=yes"`
    }
    
    // First create the structure
    entry := InsertPerson{
        Name:    "Test",
        Dtadded: time.Now(),
        Status:  1,
    }
    // Now create the query
    sqlQuery, err := gsdb.DB.Insert(entry)
    if err != nil {
        l.Error(err.Error())
        return
    }
    // Then execute the query
    
    lastInsertedID, rowsAffected, err := gsdb.DB.Execute(sqlQuery)
    if err != nil {
        l.Error(err.Error())
    }
    
    l.Info(fmt.Sprintf("Item with ID %d was inserted. %d rows were affected", lastInsertedID, rowsAffected))
    
    for i := 0; i < 100; i++ {
        _, _ = gsdb.QueryStruct[InsertPerson]("select * from Test")
        P, _ := gsdb.QuerySingleStruct[InsertPerson]("select * from Test WHERE ID =? ", lastInsertedID)
        
        if P.Id != int(lastInsertedID) {
            l.Error("ID mismatch")
        }
        
    }
    
}
