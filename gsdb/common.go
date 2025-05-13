package gsdb

import (
	"errors"
	"fmt"
	l "log/slog"
	"reflect"
	"strings"
	"time"
	"unicode"
)

// generateBuildSql creates the part of the insert SQL query which specifies which columns are to be inserted
func generateBuildSql(dbStructure any, t reflect.Type) (table string, buildSql string, err error) {
	var sb strings.Builder

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("db")
		dbStructureMap := decodeTag(tag)

		if reflect.ValueOf(dbStructure).Field(i).CanInterface() {

			if dbStructureMap["column"] == "" {
				return "", "", errors.New("no column name specified for field " + field.Name)
			}

			if dbStructureMap["table"] != "" {
				table = dbStructureMap["table"]
			}

			if dbStructureMap["omit"] != "yes" && dbStructureMap["primarykey"] != "yes" {
				// sb.WriteRune('`')
				sb.WriteString(dbStructureMap["column"])
				// sb.WriteRune('`')
				sb.WriteString(",")
			}
		}
	}

	return table, strings.TrimSuffix(sb.String(), ","), err
}

// generateValuesSql creates the part of insert SQL query that adds each entry for each structure
func generateValuesSql(dbStructure any, t reflect.Type) (string, error) {
	var sb strings.Builder
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("db")
		dbStructureMap := decodeTag(tag)

		if reflect.ValueOf(dbStructure).Field(i).CanInterface() {
			value := reflect.ValueOf(dbStructure).Field(i).Interface()

			if dbStructureMap["column"] == "" {
				return "", errors.New("no column name specified for field" + field.Type.Name())
			}

			if dbStructureMap["omit"] != "yes" && dbStructureMap["primarykey"] != "yes" {
				switch field.Type.Name() {
				case "uint", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int", "int32", "int64":
					sb.WriteString(fmt.Sprintf("%v,", value))
				case "string":
					sb.WriteString(hexRepresentation(value.(string)) + ",")
				case "float32", "float64":
					sb.WriteString(fmt.Sprintf("%v,", value))
				case "bool":
					sb.WriteString(fmt.Sprintf("%v,", value))
				case "Time":
					sb.WriteString(fmt.Sprintf("'%s',", value.(time.Time).Format("2006-01-02 15:04:05")))
				default:
					l.With("type", field.Type.Name()).With("value", value).Error("type error")
					sb.WriteString(fmt.Sprintf(`'%s',`, value.(string)))
				}
			}
		}
	}
	return fmt.Sprintf("(%s)", strings.TrimSuffix(sb.String(), ",")), nil
}

// decodeTags Turn a tag string into a map of key/value pairs
func decodeTag(tag string) map[string]string {

	lastQuote := rune(0)
	f := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)

		}
	}

	// splitting string by space but considering quoted section
	items := strings.FieldsFunc(tag, f)

	// create and fill the map
	m := make(map[string]string)
	for _, item := range items {
		x := strings.Split(item, "=")
		m[x[0]] = x[1]
	}

	// print the map
	// for k, v := range m {
	//    fmt.Printf("%s: %s\n", k, v)
	// }
	return m
}

// HexRepresentation Convert a string to a hex representation
func hexRepresentation(in string) string {
	return "X'" + fmt.Sprintf("%x", in) + "'"
	// return "'" + in + "'"
}

// getStructDetails Get the details of a struct
func getStructDetails[T any](dbFieldName string) (string, map[string]string, any) {

	var st T
	t := reflect.TypeOf(st)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("db")
		dbStructureMap := decodeTag(tag)
		// l.INFO("%d. %v (%v), tag: '%v'\n", i+1, field.Name, field.Type.Name(), tag)
		// l.SPEW(field.Type)

		if dbStructureMap["column"] == dbFieldName {
			if field.Type == reflect.TypeOf([]uint8{}) {
				return field.Name, dbStructureMap, "[]uint8"
			} else {
				return field.Name, dbStructureMap, field.Type.Name()
			}
		}
	}
	return "", map[string]string{}, ""
}
