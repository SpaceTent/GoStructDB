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

func TestInsertReadDefaultNullMatrix(t *testing.T) {
	type InsertPersonReadDefaultNull struct {
		Id      int       `db:"column=id primarykey=yes table=Test"`
		Name    string    `db:"column=name"`
		Dtadded time.Time `db:"column=dtadded readdefault=null"`
		Status  int       `db:"column=status"`
	}

	type InsertPersonNoDirective struct {
		Id      int       `db:"column=id primarykey=yes table=Test"`
		Name    string    `db:"column=name"`
		Dtadded time.Time `db:"column=dtadded"`
		Status  int       `db:"column=status"`
	}

	tests := []struct {
		name     string
		entry    any
		expected string
	}{
		{
			name: "readdefault null with zero time writes NULL",
			entry: InsertPersonReadDefaultNull{
				Name:    "Test",
				Dtadded: time.Time{},
				Status:  1,
			},
			expected: "INSERT INTO Test(name,dtadded,status) VALUES (X'54657374',NULL,1);",
		},
		{
			name: "readdefault null with non-zero time writes timestamp",
			entry: InsertPersonReadDefaultNull{
				Name:    "Test",
				Dtadded: time.Date(2025, time.December, 25, 15, 29, 25, 10, time.UTC),
				Status:  1,
			},
			expected: "INSERT INTO Test(name,dtadded,status) VALUES (X'54657374','2025-12-25 15:29:25',1);",
		},
		{
			name: "no directive with zero time writes zero timestamp",
			entry: InsertPersonNoDirective{
				Name:    "Test",
				Dtadded: time.Time{},
				Status:  1,
			},
			expected: "INSERT INTO Test(name,dtadded,status) VALUES (X'54657374','0001-01-01 00:00:00',1);",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sqlQuery, err := DB.Insert(tc.entry)
			if err != nil {
				t.Fatalf("failed to generate Insert SQL: %v", err)
			}

			if sqlQuery != tc.expected {
				t.Errorf("Expected: %s, got: %s", tc.expected, sqlQuery)
			}
		})
	}
}
