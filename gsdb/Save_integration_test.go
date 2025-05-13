package gsdb

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strconv"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

type IntegrationTestingSaveTestCase[T comparable] struct {
	name          string
	pkeyValue     T
	sqlColumnType string
}

// setUpIntegrationTest sets up the SQL Lite connection
func setUpSaveIntegrationTestConnection(t *testing.T) string {
	tempFile, err := os.CreateTemp("", "integration-test-*.db")
	assert.NoError(t, err)
	db, err := sql.Open("sqlite3", tempFile.Name())
	assert.NoError(t, err)

	New("test/test", slog.Default(), context.Background())
	DB.dbConnection = db
	DB.connected = true
	return tempFile.Name()
}

// setUpIntegrationSaveTable sets up the table for the testing
func setUpIntegrationSaveTable(t *testing.T, SQLColumnType string) {
	tearDownIntegrationSaveTable(t)
	_, err := DB.dbConnection.Exec(fmt.Sprintf(`
        CREATE TABLE Users (
            id %s PRIMARY KEY,
            name TEXT,
            status INT,
            dtadded DATETIME
        )
    `, SQLColumnType))
	assert.NoError(t, err)
}

// tearDownIntegrationSaveTable deletes the table after use
func tearDownIntegrationSaveTable(t *testing.T) {
	_, err := DB.dbConnection.Exec(`DROP TABLE IF EXISTS Users;`)
	assert.NoError(t, err)
}

// tearDownSaveIntegrationTestConnection closes the sql lite connection
func tearDownIntegrationSaveTestConnection(t *testing.T, filename string) {
	err := DB.dbConnection.Close()
	assert.NoError(t, err)
	os.Remove(filename)
}

// testIntegrationSaveTestHelper is the core tester for all kinds of types, not considering "bool" because who uses bool as a pkey value
func testIntegrationSaveTestHelper[T comparable](t *testing.T, testCases []IntegrationTestingSaveTestCase[T]) {
	// setup sql lite connection
	filename := setUpSaveIntegrationTestConnection(t)
	defer tearDownIntegrationSaveTestConnection(t, filename)

	// this is how it will appear in the table
	type IntegrationGenericStruct[V comparable] struct {
		Id     V      `db:"column=id primarykey=yes table=Users"`
		Name   string `db:"column=name"`
		Status int    `db:"column=status"`
	}

	// made this struct separate because insert ignores pkey column in the insert query and we want to set that specifically
	// in the insert query
	type IntegrationInsertGenericStruct[V comparable] struct {
		Id     V      `db:"column=id table=Users"`
		Name   string `db:"column=name"`
		Status int    `db:"column=status"`
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// create the Users table, and drop it after done
			setUpIntegrationSaveTable(t, tc.sqlColumnType)
			defer tearDownIntegrationSaveTable(t)

			value := reflect.ValueOf(tc.pkeyValue)
			// this is to separate out into insert and update modes of the DB.Save() function
			if value.IsValid() && value.IsZero() {
				entry := IntegrationInsertGenericStruct[T]{
					Id:     tc.pkeyValue,
					Name:   "Test",
					Status: 1,
				}
				// insert into table, not checking lastInsertedId because right now it only returns int64, will fail for other types
				// like strings and floats
				_, rowsAffected, err := DB.Save(entry, entry.Id)
				assert.NoError(t, err)
				assert.Equal(t, int64(1), rowsAffected)
				// check if value was successfully created in table
				result, err := QuerySingleStruct[IntegrationGenericStruct[T]]("SELECT id,name,status from Users WHERE id=?", tc.pkeyValue)
				assert.NoError(t, err)
				assert.Equal(t, tc.pkeyValue, result.Id)
				assert.Equal(t, "Test", result.Name)
				assert.Equal(t, int(1), result.Status)
				return
			}
			// update table mode of DB.Save(), we first insert into table
			var rowsAffected int64
			var err error
			// separating out because of error of unsupported values with high bit set in case of uint64
			if value.Type().Name() == "uint64" {
				_, rowsAffected, err = DB.Execute("INSERT INTO Users(id,name,status) VALUES (?,?,?)", strconv.FormatUint(value.Uint(), 10), "Test", 1)
			} else {
				_, rowsAffected, err = DB.Execute("INSERT INTO Users(id,name,status) VALUES (?,?,?)", tc.pkeyValue, "Test", 1)
			}
			assert.NoError(t, err)
			assert.Equal(t, int64(1), rowsAffected)
			// now update value
			updatedEntry := IntegrationGenericStruct[T]{
				Id:     tc.pkeyValue,
				Name:   "Test1",
				Status: 1,
			}
			_, rowsAffected, err = DB.Save(updatedEntry, updatedEntry.Id)
			assert.NoError(t, err)
			assert.Equal(t, int64(1), rowsAffected)
			// same bit set error handling, and now we check if value successfully changed in table
			var result IntegrationGenericStruct[T]
			if value.Type().Name() == "uint64" {
				result, err = QuerySingleStruct[IntegrationGenericStruct[T]]("SELECT id,name,status from Users WHERE id=?", strconv.FormatUint(value.Uint(), 10))
			} else {
				result, err = QuerySingleStruct[IntegrationGenericStruct[T]]("SELECT id,name,status from Users WHERE id=?", tc.pkeyValue)
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.pkeyValue, result.Id)
			assert.Equal(t, "Test1", result.Name)
			assert.Equal(t, int(1), result.Status)
		})
	}
}

// TestSaveIntegrationUintType tests all uint types
func TestSaveIntegrationUintType(t *testing.T) {
	testUintCases := []IntegrationTestingSaveTestCase[uint]{
		{"Uint Zero", uint(0), `INT UNSIGNED`},
		{"Uint Non-Zero", uint(42), `INT UNSIGNED`},
	}
	testIntegrationSaveTestHelper[uint](t, testUintCases)
	testUint8Cases := []IntegrationTestingSaveTestCase[uint8]{
		{"Uint8 Zero", uint8(0), `TINYINT UNSIGNED`},
		{"Uint8 Non-Zero", uint8(255), `TINYINT UNSIGNED`},
	}
	testIntegrationSaveTestHelper[uint8](t, testUint8Cases)
	testUint16Cases := []IntegrationTestingSaveTestCase[uint16]{
		{"Uint16 Zero", uint16(0), `SMALLINT UNSIGNED`},
		{"Uint16 Non-Zero", uint16(65535), `SMALLINT UNSIGNED`},
	}
	testIntegrationSaveTestHelper[uint16](t, testUint16Cases)
	testUint32Cases := []IntegrationTestingSaveTestCase[uint32]{
		{"Uint32 Zero", uint32(0), `INT UNSIGNED`},
		{"Uint32 Non-Zero", uint32(4294967295), `INT UNSIGNED`},
	}
	testIntegrationSaveTestHelper[uint32](t, testUint32Cases)
	testUint64Cases := []IntegrationTestingSaveTestCase[uint64]{
		{"Uint64 Zero", uint64(0), `BIGINT UNSIGNED`},
		{"Uint64 Non-Zero", uint64(18446744073709551615), `BIGINT UNSIGNED`},
	}
	testIntegrationSaveTestHelper[uint64](t, testUint64Cases)
}

// TestSaveIntegrationIntType tests all int types
func TestSaveIntegrationIntType(t *testing.T) {
	testIntCases := []IntegrationTestingSaveTestCase[int]{
		{"Int Zero", 0, `INT`},
		{"Int Positive", 42, `INT`},
		{"Int Negative", -42, `INT`},
	}
	testIntegrationSaveTestHelper[int](t, testIntCases)
	testInt8Cases := []IntegrationTestingSaveTestCase[int8]{
		{"Int8 Zero", int8(0), `TINYINT`},
		{"Int8 Positive", int8(127), `TINYINT`},
		{"Int8 Negative", int8(-128), `TINYINT`},
	}
	testIntegrationSaveTestHelper[int8](t, testInt8Cases)
	testInt16Cases := []IntegrationTestingSaveTestCase[int16]{
		{"Int16 Zero", int16(0), `SMALLINT`},
		{"Int16 Positive", int16(32767), `SMALLINT`},
		{"Int16 Negative", int16(-32768), `SMALLINT`},
	}
	testIntegrationSaveTestHelper[int16](t, testInt16Cases)
	testInt32Cases := []IntegrationTestingSaveTestCase[int32]{
		{"Int32 Zero", int32(0), `INT`},
		{"Int32 Positive", int32(2147483647), `INT`},
		{"Int32 Negative", int32(-2147483648), `INT`},
	}
	testIntegrationSaveTestHelper[int32](t, testInt32Cases)
	testInt64Cases := []IntegrationTestingSaveTestCase[int64]{
		{"Int64 Zero", int64(0), `BIGINT`},
		{"Int64 Positive", int64(9223372036854775807), `BIGINT`},
		{"Int64 Negative", int64(-9223372036854775808), `BIGINT`},
	}
	testIntegrationSaveTestHelper[int64](t, testInt64Cases)
}

// TestSaveIntegrationTestFloatType tests all float types
func TestSaveIntegrationTestFloatType(t *testing.T) {
	testFloat32Cases := []IntegrationTestingSaveTestCase[float32]{
		{"Float32 Zero", float32(0), `FLOAT`},
		// not checking non-zero values because of precision issues failing to find the exact value of id
	}
	testIntegrationSaveTestHelper[float32](t, testFloat32Cases)
	testFloat64Cases := []IntegrationTestingSaveTestCase[float64]{
		{"Float64 Zero", float64(0), `DOUBLE`},
		{"Float64 Positive", float64(3.14159), `DOUBLE`},
		{"Float64 Negative", float64(-3.14159), `DOUBLE`},
	}
	testIntegrationSaveTestHelper[float64](t, testFloat64Cases)
}

// TestSaveIntegrationTestStringType tests all string types
func TestSaveIntegrationTestStringType(t *testing.T) {
	testStringCases := []IntegrationTestingSaveTestCase[string]{
		// not handling empty string, because its not considered as an actual entry in a primary key value
		{"String Non-Empty", "42", `VARCHAR(255)`},
	}
	testIntegrationSaveTestHelper[string](t, testStringCases)
}
