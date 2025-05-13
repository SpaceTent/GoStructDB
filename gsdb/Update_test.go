package gsdb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type UpdatePerson[StatusType uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64 | float32 | float64 | string | bool] struct {
	Id      int        `db:"column=id primarykey=yes table=Users"`
	Name    string     `db:"column=name"`
	Dtadded time.Time  `db:"column=dtadded omit=yes"`
	Status  StatusType `db:"column=status"`
}

type UpdatePersonTime struct {
	Id      int       `db:"column=id primarykey=yes table=Users"`
	Name    string    `db:"column=name"`
	Dtadded time.Time `db:"column=dtadded"`
}

func generateUpdatePerson[StatusType uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64 | float32 | float64 | string | bool](value StatusType) UpdatePerson[StatusType] {
	return UpdatePerson[StatusType]{
		0, "Test", time.Now(), value,
	}
}

func generateUpdatePersonTime(id int) UpdatePersonTime {
	return UpdatePersonTime{
		id, "Test", time.Date(2024, time.December, 7, 15, 29, 25, 10, time.UTC),
	}
}

func testUpdateNumericalErrorValueHelper(t *testing.T, sql string, err error) {
	assert.NoError(t, err)
	assert.Equal(t, "UPDATE Users SET name=X'54657374',status=1 WHERE id=0;", sql)
}

func testUpdateBoolErrorValueHelper(t *testing.T, sql string, err error) {
	assert.NoError(t, err)
	assert.Equal(t, "UPDATE Users SET name=X'54657374',status=true WHERE id=0;", sql)
}

func testUpdateStringErrorValueHelper(t *testing.T, sql string, err error) {
	assert.NoError(t, err)
	assert.Equal(t, "UPDATE Users SET name=X'54657374',status=X'31' WHERE id=0;", sql)
}

func testUpdateTimeErrorValueHelper(t *testing.T, sql string, err error) {
	assert.NoError(t, err)
	assert.Equal(t, "UPDATE Users SET name=X'54657374',dtadded='2024-12-07 15:29:25' WHERE id=0;", sql)
}

func TestUpdate(t *testing.T) {
	New("", nil, context.Background())

	sql, err := DB.Update(generateUpdatePerson(uint(1)))
	testUpdateNumericalErrorValueHelper(t, sql, err)
	sql, err = DB.Update(generateUpdatePerson(uint8(1)))
	testUpdateNumericalErrorValueHelper(t, sql, err)
	sql, err = DB.Update(generateUpdatePerson(uint16(1)))
	testUpdateNumericalErrorValueHelper(t, sql, err)
	sql, err = DB.Update(generateUpdatePerson(uint32(1)))
	testUpdateNumericalErrorValueHelper(t, sql, err)
	sql, err = DB.Update(generateUpdatePerson(uint64(1)))
	testUpdateNumericalErrorValueHelper(t, sql, err)
	sql, err = DB.Update(generateUpdatePerson(int(1)))
	testUpdateNumericalErrorValueHelper(t, sql, err)
	sql, err = DB.Update(generateUpdatePerson(int8(1)))
	testUpdateNumericalErrorValueHelper(t, sql, err)
	sql, err = DB.Update(generateUpdatePerson(int16(1)))
	testUpdateNumericalErrorValueHelper(t, sql, err)
	sql, err = DB.Update(generateUpdatePerson(int32(1)))
	testUpdateNumericalErrorValueHelper(t, sql, err)
	sql, err = DB.Update(generateUpdatePerson(int64(1)))
	testUpdateNumericalErrorValueHelper(t, sql, err)
	sql, err = DB.Update(generateUpdatePerson(float32(1)))
	testUpdateNumericalErrorValueHelper(t, sql, err)
	sql, err = DB.Update(generateUpdatePerson(float64(1)))
	testUpdateNumericalErrorValueHelper(t, sql, err)
	sql, err = DB.Update(generateUpdatePerson(true))
	testUpdateBoolErrorValueHelper(t, sql, err)

	sql, err = DB.Update(generateUpdatePerson("1"))
	testUpdateStringErrorValueHelper(t, sql, err)

	sql, err = DB.Update(generateUpdatePersonTime(0))
	testUpdateTimeErrorValueHelper(t, sql, err)
}

func TestNoColumnNameUpdate(t *testing.T) {
	type testType1 struct {
		Id      int       `db:"column=id primarykey=yes table=Users"`
		Name    string    `db:"column="`
		Dtadded time.Time `db:"column=dtadded"`
	}
	type testType2 struct {
		Id      int       `db:"column=id primarykey=yes table=Users"`
		Name    string    `db:""`
		Dtadded time.Time `db:"column=dtadded"`
	}
	type testType3 struct {
		Id      int `db:"column=id primarykey=yes table=Users"`
		Name    string
		Dtadded time.Time `db:"column=dtadded"`
	}

	testTypeEntry1 := testType1{0, "Test", time.Now()}
	testTypeEntry2 := testType2{0, "Test", time.Now()}
	testTypeEntry3 := testType3{0, "Test", time.Now()}

	sql, err := DB.Update(testTypeEntry1)
	assert.EqualError(t, err, "no column name specified for field Name")
	assert.Empty(t, sql)

	sql, err = DB.Update(testTypeEntry2)
	assert.EqualError(t, err, "no column name specified for field Name")
	assert.Empty(t, sql)

	sql, err = DB.Update(testTypeEntry3)
	assert.EqualError(t, err, "no column name specified for field Name")
	assert.Empty(t, sql)
}

func TestNoTableNameUpdate(t *testing.T) {
	type testType1 struct {
		Id      int       `db:"column=id primarykey=yes"`
		Name    string    `db:"column=name"`
		Dtadded time.Time `db:"column=dtadded"`
	}

	testTypeEntry1 := testType1{0, "Test", time.Now()}

	sql, err := DB.Update(testTypeEntry1)
	assert.EqualError(t, err, "no table found in structure")
	assert.Empty(t, sql)
}

func TestNoFieldsUpdate(t *testing.T) {
	type testType1 struct {
	}
	type testType2 struct {
		Id      int       `db:"column=id primarykey=yes table=Users"`
		Name    string    `db:"column=name omit=yes"`
		Dtadded time.Time `db:"column=dtadded omit=yes"`
	}

	testTypeEntry1 := testType1{}
	testTypeEntry2 := testType2{0, "Test", time.Now()}

	sql, err := DB.Update(testTypeEntry1)
	assert.EqualError(t, err, "no table found in structure")
	assert.Empty(t, sql)

	sql, err = DB.Update(testTypeEntry2)
	assert.EqualError(t, err, "no non-primary key and non-omitted fields found in structure")
	assert.Empty(t, sql)
}
