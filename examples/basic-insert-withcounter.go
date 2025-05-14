package examples

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"

	MySQL "GoStructDB/gsdb"
)

func BasicInsertWithCounter() {

	godotenv.Load(".env")
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

	MySQL.DB.StartCounter("test")
	lastInsertedID, rowsAffected, err := MySQL.DB.Execute(sqlQuery)
	if err != nil {
		l.Error(err.Error())
	}

	l.Info(fmt.Sprintf("Item with ID %d was inserted. %d rows were affected", lastInsertedID, rowsAffected))

	for i := 0; i < 100; i++ {
		_, _ = MySQL.QueryStruct[InsertPerson]("select * from Test")
		_, _ = MySQL.QuerySingleStruct[InsertPerson]("select * from Test LIMIT 1")
	}

	l.With("Counter", MySQL.DB.GetCounter("test")).Info("Number of SQL Queries")

}
