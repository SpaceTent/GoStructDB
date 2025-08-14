package gsdb

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	dbConnection               *sql.DB
	DSN                        string
	Logger                     *slog.Logger
	ColumnWarnings             bool
	Lock                       sync.Mutex
	connected                  bool
	MaxDatabaseOpenConnections int
	MaxDatabaseIdleConnections int
	DatabaseIdleTimeout        time.Duration
	Ctx                        context.Context
	Counters
}

type Counters struct {
	Count map[string]int64
	Lock  sync.Mutex
}

var DB *Database
var ColumnWarnings bool = false
var ShowSQL bool = false

func New(newDSN string, L *slog.Logger, c context.Context) {

	DB = &Database{
		connected: false,
		DSN:       newDSN,
		Logger:    L,
		Ctx:       c,
	}

	DB.Counters.Count = make(map[string]int64)
}

func NewMySQL(dsn string, L *slog.Logger, c context.Context) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		L.With("Error", err, "dsn", dsn).Error("failed to open MySQL database")
	}

	DB = &Database{
		connected: true,
		dbConnection: db,
		DSN: dsn,
		Logger: L,
		Ctx: c,
	}

	DB.Counters.Count = make(map[string]int64)
}

func NewSQLite3(fileName string, L *slog.Logger, c context.Context) {

	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		L.With("Error", err, "filename", fileName).Error("unable to open/create database")
	}

	DB = &Database{
		connected:    true,
		dbConnection: db,
		DSN:          fileName,
		Logger:       L,
		Ctx:          c,
	}

	DB.Counters.Count = make(map[string]int64)
}

func getConnection() (*sql.DB, error) {

	DB.Lock.Lock()
	// check once more - in case a prev goroutine has established a connection
	if DB.connected && DB.dbConnection != nil {
		DB.Lock.Unlock()
		return DB.dbConnection, nil
	}

	if DB.DSN == "" {
		return nil, errors.New("empty database dsn")
	}

	var err error

	// attempt 3 times to connect, then give up
	for i := 0; i < 3; i++ {

		DB.dbConnection, err = sql.Open("mysql", DB.DSN)

		if err == nil {
			// Open may just validate its arguments without creating a connection to the database.
			// To verify that the data source name is valid, call Ping.
			err = DB.dbConnection.Ping()
			if err == nil {
				break // connection was fine
			}
			DB.Logger.With("attempt", i).With("error", err.Error()).Error("Unable to Ping Database")
			time.Sleep(500 * time.Millisecond) // wait a short while before trying again
			continue
		}
		time.Sleep(500 * time.Millisecond)
		DB.Logger.With("attempt", i).With("error", err.Error()).Error("Unable to Ping Database")
	}

	if err != nil {
		return nil, err
	}

	DB.dbConnection.SetMaxOpenConns(25)
	DB.dbConnection.SetMaxIdleConns(25)
	DB.dbConnection.SetConnMaxIdleTime(5 * time.Minute)
	DB.connected = true
	DB.Lock.Unlock()

	return DB.dbConnection, nil
}
