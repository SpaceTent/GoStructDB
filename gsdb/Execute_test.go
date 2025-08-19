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
	_, _, err := DB.Execute("DROP TABLE IF EXISTS test;")
	if err != nil {
		t.Fatalf("failed to execute drop table prior to tests: %v", err)
	}
	_, _, err = DB.Execute(tableCreate)
	if err != nil {
		t.Fatalf("failed to execute tableCreate SQL prior to tests: %v", err)
	}

	tests := []struct {
		name    string
		entry   TestPerson
		expectedLastID int64
	}{
		{
			name:  "Execute an Insert success with all fields",
			entry: TestPerson{Name: "Ronald McDonald", Dtadded: time.Now().UTC(), Status: 1},
			expectedLastID: 1,
		},
		{
			name:  "Execute an Insert with fields missing - current setup will populate missing fields with Go zero values",
			entry: TestPerson{Dtadded: time.Now().UTC()},
			expectedLastID: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			insertSQL, err := DB.Insert(tc.entry)
			if err != nil {
				t.Fatalf("failed to generate Insert SQL: %v", err)
			}

			lastInsertedID, rowsAffected, err := DB.Execute(insertSQL)
			if err != nil {
				t.Fatalf("failed to execute insert during test: %s: %v", tc.name, err)
			}

			var queriedName string
			err = DB.dbConnection.QueryRowContext(
				context.Background(),
				`SELECT name FROM test ORDER BY id DESC LIMIT 1`,
			).Scan(&queriedName)
			if err != nil {
				t.Fatalf("QueryRowContext scan failed in TestExecute: %v", err)
			}
			
			if tc.expectedLastID != lastInsertedID || queriedName != tc.entry.Name || rowsAffected != 1 {
				t.Fatalf("lastInsertedID does not match expected entry ID: got: %d want: %d", lastInsertedID, tc.expectedLastID)
			}
		})
	}
}
