package gsdb

import (
	"context"
	"database/sql/driver"
	"errors"
	"log/slog"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type SavePersonTime struct {
	Id      int       `db:"column=id primarykey=yes table=Users"`
	Name    string    `db:"column=name"`
	Dtadded time.Time `db:"column=dtadded omit=yes"`
	Status  int       `db:"column=status"`
}

// setupSaveTestMock sets up the mocks using sqlmock library
func setupSaveTestMock(t *testing.T, sql string, params ...driver.Value) (*sqlmock.Sqlmock, *sqlmock.ExpectedExec) {
	New("test/test", slog.Default(), context.Background())
	var err error
	var mock sqlmock.Sqlmock
	DB.dbConnection, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	DB.connected = true
	assert.NoError(t, err)
	expectedExec := mock.ExpectExec(sql)
	if len(params) > 0 {
		expectedExec.WithArgs(params...)
	}
	return &mock, expectedExec
}

// TestSaveNormalInsert tests normal insert
func TestSaveNormalInsert(t *testing.T) {
	mock, expectedExec := setupSaveTestMock(t, `INSERT INTO Users(name,status) VALUES (X'54657374',31);`)
	expectedExec.WillReturnResult(sqlmock.NewResult(1, 1))
	entry := SavePersonTime{0, "Test", time.Now(), 31}
	lastInsertedID, rowsAffected, err := DB.Save(entry, entry.Id)
	assert.NoError(t, err)
	assert.NoError(t, (*mock).ExpectationsWereMet())
	assert.Equal(t, lastInsertedID, int64(1))
	assert.Equal(t, rowsAffected, int64(1))
}

// TestSaveNormalUpdate tests normal update
func TestSaveNormalUpdate(t *testing.T) {
	mock, expectedExec := setupSaveTestMock(t, `UPDATE Users SET name=X'54657374',status=31 WHERE id=1;`)
	expectedExec.WillReturnResult(sqlmock.NewResult(0, 1))
	entry := SavePersonTime{1, "Test", time.Now(), 31}
	lastInsertedID, rowsAffected, err := DB.Save(entry, entry.Id)
	assert.NoError(t, err)
	assert.NoError(t, (*mock).ExpectationsWereMet())
	assert.Equal(t, lastInsertedID, int64(0))
	assert.Equal(t, rowsAffected, int64(1))
}

// TestSaveNoColumn tests no column field struct
func TestSaveNoColumn(t *testing.T) {
	mock, expectedExec := setupSaveTestMock(t, `UPDATE Users SET name=X'54657374',status=31 WHERE id=1;`)
	expectedExec.WillReturnResult(sqlmock.NewResult(1, 1))
	type NoColumn struct {
		Id      int       `db:"column=id primarykey=yes table=Users"`
		Name    string    `db:"column="`
		Dtadded time.Time `db:"column=dtadded omit=yes"`
		Status  int       `db:"column=status"`
	}
	entry := NoColumn{1, "Test", time.Now(), 31}
	_, _, err := DB.Save(entry, entry.Id)
	assert.EqualError(t, err, "no column name specified for field Name")
	assert.Error(t, (*mock).ExpectationsWereMet())

	_, _, err = DB.Save(entry, 0)
	assert.EqualError(t, err, "no column name specified for field Name")
	assert.Error(t, (*mock).ExpectationsWereMet())
}

// TestSaveNoPKVUpdate tests no pkv in struct with update
func TestSaveNoPKVUpdate(t *testing.T) {
	mock, expectedExec := setupSaveTestMock(t, `UPDATE Users SET name=X'54657374',status=31 WHERE id=1;`)
	expectedExec.WillReturnResult(sqlmock.NewResult(1, 1))
	entry := struct {
		Id      int       `db:"column=id table=Users"`
		Name    string    `db:"column=name"`
		Dtadded time.Time `db:"column=dtadded omit=yes"`
		Status  int       `db:"column=status"`
	}{1, "Test", time.Now(), 31}
	lastInsertedID, rowsAffected, err := DB.Save(entry, entry.Id)
	assert.EqualError(t, err, "no primary key set, unable to set a where clause")
	assert.Error(t, (*mock).ExpectationsWereMet())
	assert.Equal(t, lastInsertedID, int64(0))
	assert.Equal(t, rowsAffected, int64(0))
}

// TestSaveEmptyStruct tests empty struct
func TestSaveEmptyStruct(t *testing.T) {
	mock, expectedExec := setupSaveTestMock(t, `UPDATE Users SET name=X'54657374',status=31 WHERE id=1;`)
	expectedExec.WillReturnResult(sqlmock.NewResult(1, 1))
	entry := struct{}{}
	_, _, err := DB.Save(entry, 0)
	assert.Error(t, err)
	assert.Error(t, (*mock).ExpectationsWereMet())

	_, _, err = DB.Save(entry, 1)
	assert.Error(t, err)
	assert.Error(t, (*mock).ExpectationsWereMet())
}

type GenericEntity struct {
	Id     interface{} `db:"column=id primarykey=yes table=Users"`
	Name   string      `db:"column=name"`
	Status int         `db:"column=status"`
}

// TestSavePrimaryKeyTypes tests with zero pkv values of uint,uint8,uint16,uint32,uint64,int,int8,int16,int32,
// int64,string,float32,float64,nil,struct,time.Time
func TestSavePrimaryKeyTypes(t *testing.T) {
	testCases := []struct {
		name             string
		primaryKeyValue  interface{}
		expectedQuery    string
		expectedIsInsert bool
	}{
		// Unsigned Integers
		{"Uint Zero", uint(0), `INSERT INTO Users(name,status) VALUES (X'54657374',31);`, true},
		{"Uint Non-Zero", uint(42), `UPDATE Users SET name=X'54657374',status=31 WHERE id=42;`, false},
		{"Uint8 Zero", uint8(0), `INSERT INTO Users(name,status) VALUES (X'54657374',31);`, true},
		{"Uint8 Non-Zero", uint8(255), `UPDATE Users SET name=X'54657374',status=31 WHERE id=255;`, false},
		{"Uint16 Zero", uint16(0), `INSERT INTO Users(name,status) VALUES (X'54657374',31);`, true},
		{"Uint16 Non-Zero", uint16(65535), `UPDATE Users SET name=X'54657374',status=31 WHERE id=65535;`, false},
		{"Uint32 Zero", uint32(0), `INSERT INTO Users(name,status) VALUES (X'54657374',31);`, true},
		{"Uint32 Non-Zero", uint32(4294967295), `UPDATE Users SET name=X'54657374',status=31 WHERE id=4294967295;`, false},
		{"Uint64 Zero", uint64(0), `INSERT INTO Users(name,status) VALUES (X'54657374',31);`, true},
		{"Uint64 Non-Zero", uint64(18446744073709551615), `UPDATE Users SET name=X'54657374',status=31 WHERE id=18446744073709551615;`, false},

		// Signed Integers
		{"Int Zero", 0, `INSERT INTO Users(name,status) VALUES (X'54657374',31);`, true},
		{"Int Positive", 42, `UPDATE Users SET name=X'54657374',status=31 WHERE id=42;`, false},
		{"Int Negative", -42, `UPDATE Users SET name=X'54657374',status=31 WHERE id=-42;`, false},
		{"Int8 Zero", int8(0), `INSERT INTO Users(name,status) VALUES (X'54657374',31);`, true},
		{"Int8 Positive", int8(127), `UPDATE Users SET name=X'54657374',status=31 WHERE id=127;`, false},
		{"Int8 Negative", int8(-128), `UPDATE Users SET name=X'54657374',status=31 WHERE id=-128;`, false},
		{"Int16 Zero", int16(0), `INSERT INTO Users(name,status) VALUES (X'54657374',31);`, true},
		{"Int16 Positive", int16(32767), `UPDATE Users SET name=X'54657374',status=31 WHERE id=32767;`, false},
		{"Int16 Negative", int16(-32768), `UPDATE Users SET name=X'54657374',status=31 WHERE id=-32768;`, false},
		{"Int32 Zero", int32(0), `INSERT INTO Users(name,status) VALUES (X'54657374',31);`, true},
		{"Int32 Positive", int32(2147483647), `UPDATE Users SET name=X'54657374',status=31 WHERE id=2147483647;`, false},
		{"Int32 Negative", int32(-2147483648), `UPDATE Users SET name=X'54657374',status=31 WHERE id=-2147483648;`, false},
		{"Int64 Zero", int64(0), `INSERT INTO Users(name,status) VALUES (X'54657374',31);`, true},
		{"Int64 Positive", int64(9223372036854775807), `UPDATE Users SET name=X'54657374',status=31 WHERE id=9223372036854775807;`, false},
		{"Int64 Negative", int64(-9223372036854775808), `UPDATE Users SET name=X'54657374',status=31 WHERE id=-9223372036854775808;`, false},

		// Floating Point
		{"Float32 Zero", float32(0), `INSERT INTO Users(name,status) VALUES (X'54657374',31);`, true},
		{"Float32 Positive", float32(3.14), `UPDATE Users SET name=X'54657374',status=31 WHERE id=3.14;`, false},
		{"Float32 Negative", float32(-3.14), `UPDATE Users SET name=X'54657374',status=31 WHERE id=-3.14;`, false},
		{"Float64 Zero", float64(0), `INSERT INTO Users(name,status) VALUES (X'54657374',31);`, true},
		{"Float64 Positive", float64(3.14159), `UPDATE Users SET name=X'54657374',status=31 WHERE id=3.14159;`, false},
		{"Float64 Negative", float64(-3.14159), `UPDATE Users SET name=X'54657374',status=31 WHERE id=-3.14159;`, false},

		// String
		{"String Empty", "", `INSERT INTO Users(name,status) VALUES (X'54657374',31);`, true},
		{"String Non-Empty", "42", `UPDATE Users SET name=X'54657374',status=31 WHERE id=42;`, false},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			// setup mock
			mock, expectedExec := setupSaveTestMock(t, tc.expectedQuery)

			// prepare test entry
			entry := GenericEntity{
				Id:     tc.primaryKeyValue,
				Name:   "Test",
				Status: 31,
			}

			// setup expected results
			var expectedID int64
			switch v := tc.primaryKeyValue.(type) {
			case int:
				expectedID = int64(v)
			case int8:
				expectedID = int64(v)
			case int16:
				expectedID = int64(v)
			case int32:
				expectedID = int64(v)
			case int64:
				expectedID = v
			case uint:
				expectedID = int64(v)
			case uint8:
				expectedID = int64(v)
			case uint16:
				expectedID = int64(v)
			case uint32:
				expectedID = int64(v)
			case uint64:
				expectedID = int64(v)
			case float32:
				expectedID = int64(v)
			case float64:
				expectedID = int64(v)
			case string:
				if v == "" {
					expectedID = 0
				} else {
					expectedID, _ = strconv.ParseInt(v, 10, 64)
				}
			default:
				t.Fatalf("Unsupported type: %T", tc.primaryKeyValue)
			}
			expectedExec.WillReturnResult(sqlmock.NewResult(expectedID, 1))

			// save
			lastInsertedID, rowsAffected, err := DB.Save(entry, entry.Id)
			assert.NoError(t, err)
			assert.NoError(t, (*mock).ExpectationsWereMet())

			assert.Equal(t, expectedID, lastInsertedID, "Update should return the existing ID")
			assert.Equal(t, int64(1), rowsAffected, "Should affect exactly one row")

		})
	}
}

// TestSaveInsertError tests with db specific insert error
func TestSaveInsertError(t *testing.T) {
	mock, expectedExec := setupSaveTestMock(t, `INSERT INTO Users(name,status) VALUES (X'54657374',31);`)
	expectedExec.WillReturnError(errors.New("dummy error"))
	entry := SavePersonTime{0, "Test", time.Now(), 31}
	_, _, err := DB.Save(entry, entry.Id)
	assert.EqualError(t, err, "dummy error")
	assert.NoError(t, (*mock).ExpectationsWereMet())
}

// TestSaveUpdateError tests with db specific update error
func TestSaveUpdateError(t *testing.T) {
	mock, expectedExec := setupSaveTestMock(t, `UPDATE Users SET name=X'54657374',status=31 WHERE id=1;`)
	expectedExec.WillReturnError(errors.New("dummy error"))
	entry := SavePersonTime{1, "Test", time.Now(), 31}
	_, _, err := DB.Save(entry, entry.Id)
	assert.EqualError(t, err, "dummy error")
	assert.NoError(t, (*mock).ExpectationsWereMet())
}
