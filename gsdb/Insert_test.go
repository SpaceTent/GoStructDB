package gsdb

import (
	l "log/slog"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {

	type InsertPerson struct {
		Id      int       `db:"column=id primarykey=yes table=Test"`
		Name    string    `db:"column=name"`
		Dtadded time.Time `db:"column=dtadded"`
		Status  int       `db:"column=status"`
		Ignored int       `db:"column=ignored omit=yes"`
	}

	DtAdded := time.Date(2025, time.December, 25, 15, 29, 25, 10, time.UTC)

	// First create the structure
	entry := InsertPerson{
		Id:      12,
		Name:    "Test",
		Dtadded: DtAdded,
		Status:  1,
	}
	// Now create the query
	sqlQuery, err := DB.Insert(entry)
	if err != nil {
		l.Error(err.Error())
		return
	}
	if sqlQuery != "INSERT INTO Test(name,dtadded,status) VALUES (X'54657374','2025-12-25 15:29:25',1);" {
		t.Errorf("Expected: %s, got: %s", "INSERT INTO Test(`name`,`dtadded`,`status`) VALUES (X'54657374','2025-12-25 15:29:25',1);", sqlQuery)
	}
}
