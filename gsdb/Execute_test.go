package gsdb

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestExecute(t *testing.T) {
	textHandler := slog.NewTextHandler(os.Stdout, nil)
	l := slog.New(textHandler)
	NewSQLite3("test.db", l, context.Background())

	type TestPerson struct {
		Id      int       `db:"column=id primarykey=yes table=Test"`
		Name    string    `db:"column=name"`
		Dtadded time.Time `db:"column=dtadded"`
		Status  int       `db:"column=status"`
		Ignored int       `db:"column=ignored omit=yes"`
	}

	tableCreate := `
	CREATE TABLE IF NOT EXISTS test (
	id INTEGER PRIMARY KEY AUTOINCREMENT, 
	name TEXT NOT NULL, 
	dtadded DATETIME DEFAULT CURRENT_TIMESTAMP, 
	status INTEGER NOT NULL
	);`
	DB.Execute("DROP TABLE IF EXISTS test;")
	DB.Execute(tableCreate)

	tests := []struct {
		name    string
		entry   TestPerson
		wantErr bool
	}{
		{
			name:    "Execute an Insert success with all fields",
			entry:   TestPerson{Name: "Ronald McDonald", Dtadded: time.Now().UTC(), Status: 1},
		},
		{
			name:    "Execute an Insert with fields missing - current setup will populate missing fields with Go zero values",
			entry:   TestPerson{Dtadded: time.Now().UTC()},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			insertSQL, err := DB.Insert(tc.entry)
			if err != nil {
				t.Fatalf("Execute() error: %v, wantErr %v", err, tc.wantErr)
			}

			_, _, err = DB.Execute(insertSQL)
			if err != nil {
				t.Fatalf("Execute() error: %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}