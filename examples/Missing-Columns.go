package examples

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"

	MySQL "github.com/SpaceTent/GoStructDB/gsdb"
)

func MissingColumns() {

	godotenv.Load(".env")
	DSN := os.Getenv("DSN")
	textHandler := slog.NewTextHandler(os.Stdout, nil)
	l := slog.New(textHandler)

	MySQL.New(DSN, l, context.Background())

	type InsertPerson struct {
		Id      int       `db:"column=id primarykey=yes table=Test"`
		Name    string    `db:"column=name"`
		Dtadded time.Time `db:"column=dtadded"`
		// Status  int       `db:"column=status"`
		Ignored int `db:"column=ignored omit=yes"`
	}

	// First create the structure
	entry := InsertPerson{
		Name:    "Test",
		Dtadded: time.Now(),
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

	P, _ := MySQL.QuerySingleStruct[InsertPerson]("select * from Test WHERE ID = ?", lastInsertedID)

	MySQL.ColumnWarnings = true

	P, _ = MySQL.QuerySingleStruct[InsertPerson]("select * from Test WHERE ID = ?", lastInsertedID)

	l.With("P", P.Id, "Name", P.Name, "Ignored", P.Ignored).Info("Record")

	// There are 2 additional columns in the database.
	type InsertPerson2 struct {
		Id      int       `db:"column=id primarykey=yes table=Test2"`
		Name    string    `db:"column=name"`
		Dtadded time.Time `db:"column=dtadded"`
		Ignored int       `db:"column=ignored omit=yes"`
	}

	MySQL.ColumnWarnings = false

	P2, _ := MySQL.QueryStruct[InsertPerson2]("select * from Test")

	fmt.Println(P2)

}
