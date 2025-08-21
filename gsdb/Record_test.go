package gsdb

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestRecordUpdate(t *testing.T) {
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
	status INTEGER NOT NULL,
	ignored INTEGER
	);`

	_, err := DB.dbConnection.Exec("DROP TABLE IF EXISTS test;")
	if err != nil {
		t.Fatalf("failed to execute drop table prior to tests: %v", err)
	}

	_, err = DB.dbConnection.Exec(tableCreate)
	if err != nil {
		t.Fatalf("failed to execute tableCreate SQL prior to tests: %v", err)
	}

	tests := []struct {
		name              string
		entry             TestPerson
		initialFieldValue any
		updateColumn      string
		updateColumnValue string
	}{
		{
			name: "Successfully update a single string record field",
			entry: TestPerson{
				Name:    "Ronald McDonald",
				Dtadded: time.Now().UTC(),
				Status:  1,
			},
			updateColumn:      "name",
			updateColumnValue: "Colonel Sanders",
		},
		{
			name: "Successfully update an integer record field",
			entry: TestPerson{
				Name:    "Yosemite Sam",
				Dtadded: time.Now().UTC(),
				Status:  0,
			},
			updateColumn:      "Status",
			updateColumnValue: "1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			insertSQL, err := DB.Insert(tc.entry)
			if err != nil {
				t.Fatalf("failed to generate Insert SQL: %v", err)
			}

			lastInsertedID, _, err := DB.Execute(insertSQL)
			if err != nil {
				t.Fatalf("failed to execture insert during test: %s: %v", tc.name, err)
			}

			updateRecord := Record{
				tc.updateColumn: Field{Value: tc.updateColumnValue},
			}

			rowsAffected, err := DB.RecordUpdate(updateRecord, "test", "id", fmt.Sprintf("%d", lastInsertedID))
			if err != nil {
				t.Fatalf("failed to RecordUpdate during test: %s: %v", tc.name, err)
			}
			if rowsAffected != 1 {
				t.Fatalf("expected rowsAffected to be 1 but got %d", rowsAffected)
			}

			query := fmt.Sprintf("SELECT * FROM test WHERE id = %d", lastInsertedID)
			result, err := QuerySingleStruct[TestPerson](query)
			if err != nil {
				t.Fatalf("QuerySingleStruct failed during %s: %v", tc.name, err)
			}

			switch strings.ToLower(tc.updateColumn) {
			case "name":
				if result.Name != tc.updateColumnValue {
					t.Errorf("expected updated name to be: %s, got: %s", tc.updateColumnValue, result.Name)
				}
			case "status":
				expected, err := strconv.Atoi(tc.updateColumnValue)
				if err != nil {
					t.Fatalf("failed to convert int status to string")
				}
				if result.Status != expected {
					t.Errorf("expected updated status to be: %d, got: %d", expected, result.Status)
				}
			default:
				t.Fatalf("currently unhandled updateColumn: %s", tc.updateColumn)
			}
		})
	}
}

func TestRecordInsert(t *testing.T) {
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
	status INTEGER NOT NULL,
	ignored INTEGER
	);`

	_, err := DB.dbConnection.Exec("DROP TABLE IF EXISTS test;")
	if err != nil {
		t.Fatalf("failed to execute drop table prior to tests: %v", err)
	}

	_, err = DB.dbConnection.Exec(tableCreate)
	if err != nil {
		t.Fatalf("failed to execute tableCreate SQL prior to tests: %v", err)
	}

	tests := []struct {
		name  string
		entry Record
	}{
		{
			name: "Successfully insert a record",
			entry: Record{
				"name":    Field{Value: "Ronald McDonald"},
				"dtAdded": Field{Value: time.Now().UTC()},
				"status":  Field{Value: 1},
			},
		},
		{
			name: "Successfully insert a record with missing field - expecting go zero values to fill missing fields",
			entry: Record{
				"name":   Field{Value: "Colonel Sanders"},
				"status": Field{Value: 0},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lastInsertedID, err := DB.RecordInsert(tc.entry, "test")
			if err != nil {
				t.Fatalf("failed to insert record during test: %s: %v", tc.name, err)
			}

			result, err := QuerySingleStruct[TestPerson](`
				SELECT name, id FROM test ORDER BY id DESC LIMIT 1;
			`)
			if err != nil {
				t.Fatalf("error during QuerySingleStruct during test: %s: %v", tc.name, err)
			}

			switch {
			case result.Name != tc.entry["name"].Value:
				t.Fatalf("unexpected name value after insert - want: %s got: %s", tc.entry["name"], result.Name)

			case result.Id != int(lastInsertedID):
				t.Fatalf("last inserted ID does not match fetched record ID - want: %d got: %d", lastInsertedID, result.Id)
			}
		})
	}
}
