package gsdb

import (
	"fmt"
	l "log/slog"
	"reflect"
)

// You can't do Method Generic types in Go, so we have to use a function.

func QueryStruct[T any](sql string, parameters ...any) ([]T, error) {

	// First of all, get all the database records, ising the old Record/Field method.
	allRecords, err := DB.Query(sql, parameters...)
	if err != nil {
		return make([]T, 0), err
	}

	results := make([]T, 0)

	for i, record := range allRecords {
		var newStructRecord T

		for k, v := range record {
			// Use Reflection to set the value.

			structFieldName, dbStructureMap, structFieldType := getStructDetails[T](k)

			// fmt.Println(dbStructureMap)
			// l.Info(fmt.Sprintf("index:%d Key:%s Value:%v structFieldName:%v structFieldType:%v", i, k, "", structFieldName, structFieldType))

			switch structFieldType {
			case "int", "int8", "int16", "int32", "int64":
				// l.Info(fmt.Sprintf("Setting Int64 field: %s to %v type: %T", structFieldName, v.Value, v.Value))

				reflect.ValueOf(&newStructRecord).Elem().FieldByName(structFieldName).SetInt(v.AsInt64())

			case "uint", "uint8", "uint16", "uint32", "uint64":
				reflect.ValueOf(&newStructRecord).Elem().FieldByName(structFieldName).SetUint(v.AsUInt64())

			case "bool":
				reflect.ValueOf(&newStructRecord).Elem().FieldByName(structFieldName).SetBool(v.AsBool())

			case "float32", "float64":
				// l.Info(fmt.Sprintf("Setting flaot64 field: %s to %v", structFieldName, v.Value))
				reflect.ValueOf(&newStructRecord).Elem().FieldByName(structFieldName).SetFloat(v.AsFloat())

			case "string":
				// l.Info(fmt.Sprintf("Setting String field: %s to %v", structFieldName, v.Value))
				reflect.ValueOf(&newStructRecord).Elem().FieldByName(structFieldName).SetString(v.AsString())

			case "Time":
				// l.Info(fmt.Sprintf("Setting Time field: %s to %v", structFieldName, v.Value))

				// Does the Read Default Exist?
				// If the time is NULL or EMPTY in the database, you can have the struct return it as
				// Zero (0001-01-01 00:00:00) or the current time.
				// The first version of this library, use the current time, and that causes all sorts of
				// issues I needed to work around,  this has been implemented give you the choice of how
				// to handle it at the struct level.

				param := ""
				if dbStructureMap["readdefault"] == "now" {
					param = "now"
				}

				if dbStructureMap["readdefault"] == "zero" {
					param = "zero"
				}

				reflect.ValueOf(&newStructRecord).Elem().FieldByName(structFieldName).Set(reflect.ValueOf(v.AsDate(param)))

				// Add Blob Support.
			case "[]uint8":
				l.Info(fmt.Sprintf("Setting Blob field: %s to %v", structFieldName, v.Value))
				reflect.ValueOf(&newStructRecord).Elem().FieldByName(structFieldName).Set(reflect.ValueOf(v.AsByte()))

			default:
				if ColumnWarnings {
					l.With("col", k).With("index", i).With("structFieldName", structFieldName).With("structFieldType", structFieldType).Warn("Unknown type")
				}
			}
		}

		results = append(results, newStructRecord)
	}
	return results, nil
}

// You can't do Method Generic types in Go, so we have to use a function.

func QuerySingleStruct[T any](sql string, parameters ...any) (T, error) {

	var SingleResult T

	results, err := QueryStruct[T](sql, parameters...)
	if err != nil {
		return SingleResult, err
	}
	if len(results) == 0 {
		return SingleResult, nil
	}
	return results[0], nil
}
