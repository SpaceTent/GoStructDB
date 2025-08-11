package gsdb

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"
)

type TestSuite struct {
	DB         *Database
	RawDB      *sql.DB
	Tables     []string  // keep a track of any created tables for cleanup
	OriginalDB *Database // store original DB state for any restoration
	LogBuffer  *bytes.Buffer
}

// TestPerson represents a dummy type for CRUD testing
type TestPerson struct {
	ID        int       `db:"column=id primarykey=yes table=people"`
	Name      string    `db:"column=name"`
	Email     string    `db:"column=email"`
	Age       int       `db:"column=age"`
	CreatedAt time.Time `db:"column=created_at"`
	UpdatedAt time.Time `db:"column=updated_at"`
	Active    bool      `db:"column=active"`
}

type TestProduct struct {
	ID          int       `db:"column=id primarykey=yes table=products"`
	Name        string    `db:"column=name"`
	Price       float64   `db:"column=price"`       // let's just assume any currency for now
	Description *string   `db:"column=description"` // test nullable fields
	CreatedAt   time.Time `db:"column=created_at"`
	UpdatedAt   time.Time `db:"column=updated_at"`
	EOL         bool      `db:"column=end_of_life"`
}

// New spins up a test DB and returns a TestSuite
func NewSuite(t *testing.T) *TestSuite {
	originalDB := DB

	// create a log buffer for the suite - print logs on test failure
	logBuffer := &bytes.Buffer{}
	logger := slog.New(slog.NewTextHandler(logBuffer, nil))
	ctx := context.Background()

	// this could probably be better
	t.Cleanup(func() {
		if t.Failed() {
			fmt.Println(os.Stderr, "=== failed test logs ===")
			fmt.Print(os.Stderr, logBuffer.String())
		}
	})

	// in-memory SQLite for tests
	NewSQLite3(":memory", logger, ctx)

	if !DB.connected || DB.dbConnection == nil {
		t.Fatalf("Failed to create test DB connection")
	}

	suite := &TestSuite{
		DB:         DB,
		RawDB:      DB.dbConnection,
		Tables:     make([]string, 0),
		OriginalDB: originalDB, // storing this for clean up - not sure if necessary at the moment
	}

	suite.createTestTables(t)

	return suite
}

// IsConnected checks if the DB connection is alive
func (s *TestSuite) IsConnected() bool {
	if s.RawDB == nil {
		return false
	}

	err := s.RawDB.Ping()
	return err == nil
}

// Clear removes any data from test tables
func (s *TestSuite) Clear() error {
	for _, table := range s.Tables {
		_, err := s.RawDB.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			return fmt.Errorf("failed to clear test table %s: %v", table, err)
		}
	}
	return nil
}

// TearDown cleans up the test DB and restores the original DB state
func (s *TestSuite) TearDown() {
	if s.RawDB != nil {
		s.RawDB.Close()
	}

	// restore DB to original state
	DB = s.OriginalDB
}

func (s *TestSuite) createTestTables(t *testing.T) {
	// People table
	peopleSQL := `
		CREATE TABLE people (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT,
		age INTEGER,
		created_at DATETIME,
		updated_at DATETIME,
		active BOOLEAN DEFAULT TRUE
	)`

	// Products test table
	productsSQL := `
		CREATE TABLE products(
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		price REAL,
		description TEXT,
		created_at DATETIME,
		updated_at DATETIME,
		end_of_life BOOLEAN DEFAULT FALSE,
	)`

	tables := map[string]string{
		"people":   peopleSQL,
		"products": productsSQL,
	}

	for tableName, sql := range tables {
		_, err := s.RawDB.Exec(sql)
		if err != nil {
			t.Fatalf("Failed to create test table %s: %v", tableName, err)
		}
		s.Tables = append(s.Tables, tableName)
	}
}

/* TODO:

- Funcs
	- CountRows?
	- TableExists
	- ExecuteQuery - execute a raw SQL query
	- QueryRow - query single
	- SeedTestData

- Pretty print logs with time stamps and log levels
*/
