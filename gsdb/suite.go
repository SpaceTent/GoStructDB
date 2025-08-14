package gsdb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"testing"
)

type Suite struct {
	DB         *Database
	Tables     []string  // keep track of any created tables for cleanup
	OriginalDB *Database // store original DB state for any restoration
	Logger     *slog.Logger
}

// New spins up a test DB and returns a test suite
func NewSuite(t *testing.T, driver, dsn string) *Suite {
	ctx := context.Background()
	originalDB := DB

	logBuffer := &bytes.Buffer{}
	logger := slog.New(slog.NewTextHandler(logBuffer, nil))

	t.Cleanup(func() {
		if t.Failed() {
			fmt.Println("=== failed test logs ===")
			fmt.Print(logBuffer.String())
		}
	})

	// check driver type - potential to add others in the future e.g. Postgres
	switch driver {
	case "sqlite3":
		NewSQLite3(dsn, logger, ctx)
	case "mysql":
		NewMySQL(dsn, logger, ctx)
	default:
		t.Fatalf("unsupported driver: %s", driver)
	}

	return &Suite{
		DB:         DB,
		Logger:     logger,
		Tables:     make([]string, 0),
		OriginalDB: originalDB,
	}
}

func (s *Suite) IsConnected(db *Database) bool {
	if db == nil || !db.connected || db.dbConnection == nil {
		return false
	}
	if err := db.dbConnection.Ping(); err != nil {
		s.DB.Logger.With("error", err.Error()).Warn("test db ping failed")
		return false
	}

	return true
}

// TableExists checks if a table exists in the DB
func (s *Suite) TableExists(table string) (bool, error) {
	if !s.IsConnected(s.DB) {
		return false, errors.New("test database not connected")
	}

	// extract driver from DSN for now
	// rather than adding a Driver field to the Database struct in New.go
	var driver string
	if strings.HasSuffix(s.DB.DSN, ".db") || strings.HasPrefix(s.DB.DSN, "file:") || strings.HasPrefix(s.DB.DSN, ":memory") {
		driver = "sqlite3"
	} else {
		driver = "mysql"
	}

	var query string

	// potentially move this to a helper in future if more DB support is added
	switch driver {
	case "sqlite3":
		query = `SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`

	case "mysql":
		query = `SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?`

	default:
		return false, fmt.Errorf("unsupported driver in TableExists")
	}

	var exists int
	err := s.DB.dbConnection.QueryRow(query, table).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}

// CountRows returns the number of rows from a given table
func (s *Suite) CountRows(table string) (int64, error) {
	if !s.IsConnected(s.DB) {
		return 0, errors.New("test db not connected: CountRows")
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s;", table)
	var count int64
	err := s.DB.dbConnection.QueryRowContext(context.Background(), query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// QueryOne executes a query with the expectation that a single row is returned
func (s *Suite) QueryOne(query string, result ...interface{}) error {
	if !s.IsConnected(s.DB) {
		return errors.New("test db not connected: QueryOne")
	}
	return s.DB.dbConnection.QueryRowContext(context.Background(), query).Scan(result...)
}

// Query executes a query that is expected to return one or more rows
func (s *Suite) Query(query string, args ...interface{}) ([]map[string]interface{}, error) {
	if !s.IsConnected(s.DB) {
		return nil, errors.New("test db not connected: Query")
	}

	rows, err := s.DB.dbConnection.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	results := []map[string]interface{}{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})
		for i, col := range columns {
			rowMap[col] = values[i]
		}

		results = append(results, rowMap)
	}

	return results, nil
}

func (s *Suite) CreateTestTables() error {
	if !s.IsConnected(s.DB) {
		return errors.New("test db is not connected: CreateTestTables")
	}

	var createStatement string

	// check if driver is sqlite
	if strings.HasSuffix(s.DB.DSN, ".db") || strings.HasPrefix(s.DB.DSN, "file:") || strings.HasPrefix(s.DB.DSN, ":memory:") {
		// sqlite
		createStatement = `
		CREATE TABLE IF NOT EXISTS people (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT,
			active INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`
	} else {
		// mysql
		createStatement = `
		CREATE TABLE IF NOT EXISTS people (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255),
			active TINYINT(1) NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		);`
	}

	_, err := s.DB.dbConnection.ExecContext(context.Background(), createStatement)
	if err != nil {
		return fmt.Errorf("failed to create test table: %v", err)
	}

	s.Tables = append(s.Tables, "people")
	return nil
}

func (s *Suite) Clear() error {
	if !s.IsConnected(s.DB) {
		return errors.New("test db is not connected: Clear")
	}

	for _, table := range s.Tables {
		_, err := s.DB.dbConnection.ExecContext(context.Background(), fmt.Sprintf("DELETE FROM %s;", table))
		if err != nil {
			return fmt.Errorf("failed to clear test table: %s: %v", table, err)
		}
	}

	return nil
}

func (s *Suite) TearDown() {
	if s.DB != nil && s.DB.dbConnection != nil {
		s.DB.dbConnection.Close()
	}

	DB = s.OriginalDB
}
