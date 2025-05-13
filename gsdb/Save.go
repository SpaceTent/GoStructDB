package gsdb

import (
	"errors"
	"reflect"
)

// Save takes in a structure and if the primary key value is set to a non-zero value, then it will update the object
// else it will insert the object into the table (taking in a primary key to reduce reflection overhead)
func (db *Database) Save(dbStructure any, primaryKeyValue any) (lastInsertedID, rowsAffected int64, err error) {
	pkvValue := reflect.ValueOf(primaryKeyValue) // pkv => Primary Key Value
	if !pkvValue.IsValid() {
		return 0, 0, errors.New("invalid primary key value")
	}
	var sql string
	if pkvValue.IsZero() {
		sql, err = DB.Insert(dbStructure)
		if err != nil {
			return 0, 0, err
		}
	} else {
		sql, err = DB.Update(dbStructure)
		if err != nil {
			return 0, 0, err
		}
	}
	return DB.Execute(sql)
}
