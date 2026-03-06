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

func SaveReadDefaultNull() {

	godotenv.Load(".env")
	DSN := os.Getenv("DSN")
	textHandler := slog.NewTextHandler(os.Stdout, nil)
	l := slog.New(textHandler)

	// MySQL.New(DSN, l, context.Background())
	MySQL.NewSQLite3(DSN, l, context.Background())

	type SavePerson struct {
		Id      int       `db:"column=id primarykey=yes table=Test"`
		Name    string    `db:"column=name"`
		Dtadded time.Time `db:"column=dtadded readdefault=null"`
		Status  int       `db:"column=status"`
	}

	entry := SavePerson{
		Name:    "Save Null",
		Dtadded: time.Time{},
		Status:  1,
	}

	lastInsertedID, rowsAffected, err := MySQL.DB.Save(entry, entry.Id)
	if err != nil {
		l.Error(err.Error())
		return
	}

	loaded, err := MySQL.QuerySingleStruct[SavePerson]("SELECT * FROM test WHERE id = ?", lastInsertedID)
	if err != nil {
		l.Error("Query error: " + err.Error())
	} else {
		l.With("Date Should be NULL/zero", loaded.Dtadded).Info("Date")
	}

	loaded.Status = 2
	loaded.Dtadded = time.Now()

	_, rowsAffectedUpdate, err := MySQL.DB.Save(loaded, loaded.Id)
	if err != nil {
		l.Error(err.Error())
	}

	l.Info(fmt.Sprintf("Save insert rows: %d, save update rows: %d", rowsAffected, rowsAffectedUpdate))
}
