package gsdb

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func (db *Database) Update(dbStructure any) (string, error) {

	t := reflect.TypeOf(dbStructure)
	UpdateTable := ""
	buildsql := ""
	UpdateColumn := ""
	UpdateValue := ""

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("db")
		dbStructureMap := decodeTag(tag)

		if reflect.ValueOf(dbStructure).Field(i).CanInterface() {
			value := reflect.ValueOf(dbStructure).Field(i).Interface()
			// l.INFO("%d. Value='%v'  %v (%v), tag: '%v'\n", i+1, value, field.Name, field.Type.Name(), tag)

			// TODO: Need to look at way for this to happen and not though an error
			if dbStructureMap["column"] == "" {
				return "", errors.New("no column name specified for field " + field.Name)
			}

			if dbStructureMap["primarykey"] == "yes" {
				// l.INFO("Primary Key Found: %s", dbStructureMap["table"])
				UpdateColumn = dbStructureMap["column"]
				UpdateValue = fmt.Sprintf("%v", value)
			}

			if dbStructureMap["table"] != "" {
				UpdateTable = dbStructureMap["table"]
			}

			if dbStructureMap["omit"] != "yes" && dbStructureMap["primarykey"] != "yes" {
				buildsql = buildsql + dbStructureMap["column"] + "="

				switch field.Type.Name() {
				case "uint", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int", "int32", "int64":
					buildsql = buildsql + fmt.Sprintf("%v", value) + ","
				case "string":
					buildsql = buildsql + hexRepresentation(value.(string)) + ","
				case "float32", "float64":
					buildsql = buildsql + fmt.Sprintf("%v", value) + ","
				case "bool":
					buildsql = buildsql + fmt.Sprintf("%v", value) + ","
				case "Time":
					buildsql = buildsql + fmt.Sprintf("'%s'", value.(time.Time).Format("2006-01-02 15:04:05")) + ","
				default:
					db.Logger.With("type", field.Type.Name()).With("value", value).Error("type error")
					buildsql = buildsql + "'" + value.(string) + "',"
				}
			}
		}
	}
	// Get Rid of Trailing Comma

	if UpdateTable == "" {
		return "", fmt.Errorf("no table found in structure")
	}

	if buildsql == "" {
		return "", fmt.Errorf("no non-primary key and non-omitted fields found in structure")
	}

	if UpdateColumn == "" {
		return "", fmt.Errorf("no primary key set, unable to set a where clause")
	}

	buildsql = strings.TrimSuffix(buildsql, ",")
	SQL := "UPDATE " + UpdateTable + " SET " + buildsql + " WHERE " + UpdateColumn + "=" + UpdateValue + ";"

	return SQL, nil
}
