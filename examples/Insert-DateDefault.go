package examples

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"

	// MySQL "github.com/SpaceTent/db/mysql" <--- old library
	MySQL "github.com/SpaceTent/GoStructDB/gsdb"
)

func InsertDateDefault() {

	godotenv.Load(".env")
	DSN := os.Getenv("DSN")
	ctx := context.Background()

	textHandler := slog.NewTextHandler(os.Stdout, nil)
	l := slog.New(textHandler)

	// MySQL.New(DSN, l, context.Background())
	MySQL.New(DSN, l, context.Background())

	type InsertPersonNOW struct {
		Id      int       `db:"column=id primarykey=yes table=Test"`
		Name    string    `db:"column=name"`
		Dtadded time.Time `db:"column=dtadded readdefault=now"`
		Status  int       `db:"column=status"`
		Ignored int       `db:"column=ignored omit=yes"`
	}

	type InsertPersonZERO struct {
		Id      int       `db:"column=id primarykey=yes table=Test"`
		Name    string    `db:"column=name"`
		Dtadded time.Time `db:"column=dtadded readdefault=zero"`
		Status  int       `db:"column=status"`
		Ignored int       `db:"column=ignored omit=yes"`
	}

	type InsertPersonDefault struct {
		Id      int       `db:"column=id primarykey=yes table=Test"`
		Name    string    `db:"column=name"`
		Dtadded time.Time `db:"column=dtadded"`
		Status  int       `db:"column=status"`
		Ignored int       `db:"column=ignored omit=yes"`
	}

	sqlQuery := "INSERT INTO Test(name) VALUES (X'54657374');"

	// Then execute the query
	lastInsertedID, rowsAffected, err := MySQL.DB.Execute(sqlQuery)
	if err != nil {
		l.Error(err.Error())
	}

	test, err2 := MySQL.QuerySingleStruct[InsertPersonNOW]("SELECT * FROM test WHERE id = ?", lastInsertedID)
	if err2 != nil {
		l.Error("Query error: " + err2.Error())
	} else {
		l.With("Date Should be NOW", test.Dtadded).Info("Date")
	}

	test2, err3 := MySQL.QuerySingleStruct[InsertPersonZERO]("SELECT * FROM test WHERE id = ?", lastInsertedID)
	if err3 != nil {
		l.Error("Query error: " + err3.Error())
	} else {
		l.With("Date Should be ZERO", test2.Dtadded).Info("Date")
	}

	test3, err4 := MySQL.QuerySingleStruct[InsertPersonDefault]("SELECT * FROM test WHERE id = ?", lastInsertedID)
	if err4 != nil {
		l.Error("Query error: " + err4.Error())
	} else {
		l.With("Date Should be DEFAULT", test3.Dtadded).Info("Date")
	}

	l.Info(fmt.Sprintf("Item with ID %d was inserted. %d rows were affected", lastInsertedID, rowsAffected))

	ctx.Done()
}
