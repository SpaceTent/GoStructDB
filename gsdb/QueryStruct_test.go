package gsdb

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestQuerySingleStruct(t *testing.T) {
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
		name           string
		entry          TestPerson
		expectedResult TestPerson
	}{
		{
			name: "Successfully return a complete struct (all fields inserted)",
			entry: TestPerson{
				Name:    "Ronald McDonald",
				Dtadded: time.Date(2012, 10, 31, 15, 50, 13, 0, time.UTC),
				Status:  1,
				Ignored: 0,
			},
			expectedResult: TestPerson{
				Id: 1,
				Name:    "Ronald McDonald",
				Dtadded: time.Date(2012, 10, 31, 15, 50, 13, 0, time.UTC),
				Status:  1,
				Ignored: 0,
			},
		},
		{
			name: "Successfully return a struct with some missing fields populated by Go zero values",
			entry: TestPerson{
				Name:   "Colonel Sanders",
				Dtadded: time.Date(2012, 10, 31, 15, 50, 13, 0, time.UTC),
			},
			expectedResult: TestPerson{
				Id: 2,
				Name:    "Colonel Sanders",
				Dtadded: time.Date(2012, 10, 31, 15, 50, 13, 0, time.UTC),
				Status:  0,
				Ignored: 0,
			},
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

			query := fmt.Sprintf("SELECT * FROM test WHERE id = %d", lastInsertedID)
			result, err := QuerySingleStruct[TestPerson](query)
			if err != nil {
				t.Fatalf("QuerySingleStruct failed during %s: %v", tc.name, err)
			}

			if !reflect.DeepEqual(result, tc.expectedResult) {
				t.Errorf("mismatch in queried result and expected result:\n got  %+v\n want %+v", result, tc.expectedResult)
			}
		})
	}
}

func TestQueryStruct(t *testing.T) {
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
		name           string
		entries        []TestPerson
		expectedResult []TestPerson
	}{
		{
			name: "Successfully return multiple structs",
			entries: []TestPerson{
				{
					Name:    "Ronald McDonald",
					Dtadded: time.Date(2012, 10, 31, 15, 50, 13, 0, time.UTC),
					Status:  1,
					Ignored: 0,
				},
				{
					Name:    "Colonel Sanders",
					Dtadded: time.Date(2012, 10, 31, 15, 50, 13, 0, time.UTC),
				},
			},
			expectedResult: []TestPerson{
				{
					Id:      1,
					Name:    "Ronald McDonald",
					Dtadded: time.Date(2012, 10, 31, 15, 50, 13, 0, time.UTC),
					Status:  1,
					Ignored: 0,
				},
				{
					Id:      2,
					Name:    "Colonel Sanders",
					Dtadded: time.Date(2012, 10, 31, 15, 50, 13, 0, time.UTC),
					Status:  0,
					Ignored: 0,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for _, entry := range tc.entries {
				insertSQL, err := DB.Insert(entry)
				if err != nil {
					t.Fatalf("failed to generate Insert SQL: %v", err)
				}

				_, _, err = DB.Execute(insertSQL)
				if err != nil {
					t.Fatalf("failed to execute insert during test: %s: %v", tc.name, err)
				}
			}

			query := "SELECT * FROM test ORDER BY id ASC;"
			results, err := QueryStruct[TestPerson](query)
			if err != nil {
				t.Fatalf("QueryStruct failed during %s: %v", tc.name, err)
			}

			if !reflect.DeepEqual(results, tc.expectedResult) {
				t.Errorf("mismatch in queried results and expected results:\n got  %+v\n want %+v", results, tc.expectedResult)
			}
		})
	}
}