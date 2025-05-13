package gsdb

import (
	"fmt"
	l "log/slog"
	"strconv"
	"time"
)

// The Database returns a map of []Records and each Record is a map of Fields.
// This provides a way to get the Field from the Record.

type Field struct {
	Value any
}

func (F Field) AsString() string {

	if F.Value == nil {
		return ""
	}

	switch v := F.Value.(type) {
	case int32, int64:
		// Convert to a String
		return strconv.FormatInt(int64(F.Value.(int64)), 10)
	case float64:
		// fmt.Printf("Float64: %v\n", val)
		// Arry of Bytes.
	case []uint8:
		b, _ := F.Value.([]byte)
		return string(b)
	case string:
		return F.Value.(string)

	default:
		l.Error("Can not convert type: '" + fmt.Sprintf("%T", v) + "' to a String")
	}

	return F.Value.(string)
}

func (F Field) AsFloat() float64 {

	if F.Value == nil {
		return 0
	}

	switch v := F.Value.(type) {
	case int:
		return float64(F.Value.(int))
	case int8:
		return float64(F.Value.(int8))
	case int16:
		return float64(F.Value.(int16))
	case int32:
		return float64(F.Value.(int32))
	case int64:
		return float64(F.Value.(int64))
	case uint:
		return float64(F.Value.(uint))
	case uint8:
		return float64(F.Value.(uint8))
	case uint16:
		return float64(F.Value.(uint16))
	case uint32:
		return float64(F.Value.(uint32))
	case uint64:
		return float64(F.Value.(uint64))
	case float32:
		return float64(F.Value.(float32))
	case float64:
		return float64(F.Value.(float64))
	case string:
		floatVal, err := strconv.ParseFloat(F.Value.(string), 64)
		if err != nil {
			l.With("err", err.Error()).Error(fmt.Sprintf("Cannot convert %T to float64", v))
			return 0
		}
		return floatVal
	default:
		l.Error("Can not convert type: '" + fmt.Sprintf("%T", v) + "' to a float64")
	}

	return float64(F.Value.(float64))
}

func (F Field) AsDate(d string) time.Time {

	// https://github.com/go-sql-driver/mysql#timetime-support
	// This assumes this is OFF.

	if F.Value == nil {
		if d != "" {

			if d == "now" {
				return time.Now()
			}
			if d == "zero" {
				return time.Time{}
			}

			out, _ := time.Parse("2006-01-02 15:04:05", d)
			return out
		} else {
			return time.Now()
		}
	}

	switch v := F.Value.(type) {
	case time.Time:
		return F.Value.(time.Time)
	case string:
		t, _ := time.Parse("2006-01-02 15:04:05", F.Value.(string))
		return t
	default:
		l.Error("Can not convert type: '" + fmt.Sprintf("%T", v) + "' to a Date")
		return F.Value.(time.Time)
	}

}

func (F Field) AsDateEpoch() int64 {

	if F.Value == nil {
		return 0
	}
	t, _ := time.Parse("2006-01-02 15:04:05", F.Value.(string))

	return t.Unix()
}

func (F Field) AsInt() int {

	// TODO:  If there is a NULL in the database.  the Interface is nil.

	if F.Value == nil {
		return 0
	}

	// This code is needed on each of the fields for flexiblity.  If you need to gt a Field from the database and have in the
	// code as a different Type.  Most of the time isn't going to be needed. This is (DEFAULT) conversion

	switch v := F.Value.(type) {
	case int:
		// Interface is a Int, so just so the conversion. (DEFAULT)
		return F.Value.(int)

	case int32:
		return int(F.Value.(int32))
	case int64:
		return int(F.Value.(int64))
	case float64:
		return int(F.Value.(float64))
	case float32:
		return int(F.Value.(float32))
	case []uint8:
		// Interface Value is an array of Characters Convert a string, then an Int
		b, _ := F.Value.([]byte)
		i, _ := strconv.Atoi(string(b))
		return i
	case string:
		i, _ := strconv.Atoi(F.Value.(string))
		return i

	default:
		l.Error("Can not convert type: '" + fmt.Sprintf("%T", v) + "' to an int")
	}

	return 0
}

func (F Field) AsInt64() int64 {

	// TODO:  If there is a NULL in the database.  the Interface is nil.

	if F.Value == nil {
		return 0
	}

	// This code is needed on each of the fields for flexiblity.  If you need to gt a Field from the database and have in the
	// code as a different Type.  Most of the time isn't going to be needed. This is (DEFAULT) conversion

	switch v := F.Value.(type) {
	case bool:
		if F.Value.(bool) {
			return int64(1)
		}
		return int64(0)
	case uint:
		return int64(F.Value.(uint))
	case uint8:
		return int64(F.Value.(uint8))
	case uint16:
		return int64(F.Value.(uint16))
	case uint32:
		return int64(F.Value.(uint32))
	case uint64:
		return int64(F.Value.(uint64))
	case int:
		// Interface is a Int, so just so the conversion. (DEFAULT)
		return int64(F.Value.(int))
	case int8:
		return int64(F.Value.(int8))
	case int16:
		return int64(F.Value.(int16))
	case int32:
		return int64(F.Value.(int32))
	case int64:
		return F.Value.(int64)
	case float64:
		return int64(F.Value.(float64))
	case float32:
		return int64(F.Value.(float32))
	case string:
		i, _ := strconv.Atoi(F.Value.(string))
		return int64(i)

	default:
		l.Error("Can not convert type %T %v %v", v, v, F.Value)
	}

	return 0
}

func (F Field) AsUInt64() uint64 {
	if F.Value == nil {
		return 0
	}

	// This code is needed on each of the fields for flexiblity.  If you need to gt a Field from the database and have in the
	// code as a different Type.  Most of the time isn't going to be needed. This is (DEFAULT) conversion

	switch v := F.Value.(type) {
	case bool:
		if F.Value.(bool) {
			return uint64(1)
		}
		return uint64(0)
	case uint:
		return uint64(F.Value.(uint))
	case uint8:
		return uint64(F.Value.(uint8))
	case uint16:
		return uint64(F.Value.(uint16))
	case uint32:
		return uint64(F.Value.(uint32))
	case uint64:
		return F.Value.(uint64)
	case int:
		// Interface is a Int, so just so the conversion. (DEFAULT)
		return uint64(F.Value.(int))
	case int8:
		return uint64(F.Value.(int8))
	case int16:
		return uint64(F.Value.(int16))
	case int32:
		return uint64(F.Value.(int32))
	case int64:
		return uint64(F.Value.(int64))
	case float64:
		return uint64(F.Value.(float64))
	case float32:
		return uint64(F.Value.(float32))
	case string:
		i, _ := strconv.Atoi(F.Value.(string))
		return uint64(i)

	default:
		l.Error("Can not convert type %T %v %v", v, v, F.Value)
	}

	return 0
}

func (F Field) AsBool() bool {
	if F.Value == nil {
		return false
	}

	switch v := F.Value.(type) {
	case bool:
		return F.Value.(bool)
	case uint:
		return convToBool(F.Value.(uint))
	case uint8:
		return convToBool(F.Value.(uint8))
	case uint16:
		return convToBool(F.Value.(uint16))
	case uint32:
		return convToBool(F.Value.(uint32))
	case uint64:
		return convToBool(F.Value.(uint64))
	case int:
		// Interface is a Int, so just so the conversion. (DEFAULT)
		return convToBool(F.Value.(int))
	case int8:
		return convToBool(F.Value.(int8))
	case int16:
		return convToBool(F.Value.(int16))
	case int32:
		return convToBool(F.Value.(int32))
	case int64:
		return convToBool(F.Value.(int64))
	case float64:
		return convToBool(F.Value.(float64))
	case float32:
		return convToBool(F.Value.(float32))
	case string:
		return len(F.Value.(string)) > 0
	default:
		l.Error("Can not convert type: '" + fmt.Sprintf("%T", v) + "' to a bool")
	}
	return false
}

func (F Field) AsByte() []byte {

	if F.Value == nil {
		return []byte{}
	}

	switch v := F.Value.(type) {

	case []uint8:
		// Interface Value is an array of Characters Convert a string, then an Int
		b, _ := F.Value.([]byte)
		return b

	case string:
		return []byte(F.Value.(string))

	default:
		l.Error("Can not convert type: '" + fmt.Sprintf("%T", v) + "' to a Bytes")
	}
	return []byte{}
}

func convToBool[T uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64 | float32 | float64](value T) bool {
	switch any(value).(type) {
	case uint, uint8, uint16, uint32, uint64:
		return value != 0
	case int, int8, int16, int32, int64:
		return value != 0
	case float32, float64:
		return value != 0.0
	default:
		return false
	}
}
