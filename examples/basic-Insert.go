package examples

import (
    "context"
    "fmt"
    "log/slog"
    "os"
    "time"
    
    "github.com/joho/godotenv"
    _ "github.com/joho/godotenv/autoload"
    
    MySQL "GoStructDB/gsdb"
)

func BasicInsert() {
    
    err := godotenv.Load("examples/.env")
    DSN := os.Getenv("DSN")
    
    textHandler := slog.NewTextHandler(os.Stdout, nil)
    l := slog.New(textHandler)
    
    MySQL.New(DSN, l, context.Background())
    
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
    sqlQuery, err := MySQL.DB.Insert(entry)
    if err != nil {
        l.Error(err.Error())
        return
    }
    // Then execute the query
    lastInsertedID, rowsAffected, err := MySQL.DB.Execute(sqlQuery)
    if err != nil {
        l.Error(err.Error())
    }
    l.Info(fmt.Sprintf("Item with ID %d was inserted. %d rows were affected", lastInsertedID, rowsAffected))
}
