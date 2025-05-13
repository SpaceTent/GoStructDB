package gsdb

import (
	"fmt"
	"strings"
	"time"
)

type Record map[string]Field

func (db *Database) RecordUpdate(RecordToUpdate Record, UpdateTable string, UpdateColumn string, UpdateColumnValue string) (int64, error) {

	// Build an SQL Statement Based on the Record.
	buildsql := "UPDATE " + UpdateTable + " SET "

	for key, F := range RecordToUpdate {
		buildsql = buildsql + key + " = "

		switch v := F.Value.(type) {
		case int, int32, int64:
			buildsql = buildsql + fmt.Sprintf("%v", F.Value) + ","
		case float64:
			buildsql = buildsql + fmt.Sprintf("%v", F.Value) + ","
		case string:
			buildsql = buildsql + hexRepresentation(F.Value.(string)) + ","
		case time.Time:
			buildsql = buildsql + fmt.Sprintf("'%s'", F.Value.(time.Time).Format("2006-01-02 15:04:05")) + ","
		default:
			DB.Logger.Error(fmt.Sprintf("%v is unknown", v))
			buildsql = buildsql + "'" + F.Value.(string) + "',"
		}

	}
	buildsql = strings.TrimSuffix(buildsql, ",")
	buildsql = buildsql + " WHERE " + UpdateColumn + " = " + UpdateColumnValue

	_, RowsAffected, err := DB.Execute(buildsql)
	if err != nil {
		return RowsAffected, err
	}
	return RowsAffected, nil
}

func (db *Database) RecordInsert(RecordToInsert Record, InsertTable string) (int64, error) {

	// Build an SQL Statement Based on the Record.
	buildsql := "INSERT INTO " + InsertTable + "("
	endsql := ""

	for key, F := range RecordToInsert {
		buildsql = buildsql + key + ","

		switch v := F.Value.(type) {
		case int, int32, int64:
			endsql = endsql + fmt.Sprintf("%v", F.AsInt()) + ","
		case string:
			endsql = endsql + hexRepresentation(F.Value.(string)) + ","
		case float64:
			endsql = endsql + fmt.Sprintf("%v", F.Value) + ","
		case time.Time:
			endsql = endsql + fmt.Sprintf("'%s'", F.Value.(time.Time).Format("2006-01-02 15:04:05")) + ","
		default:
			DB.Logger.Error(fmt.Sprintf("%v is unknown", v))
			endsql = endsql + "'" + F.Value.(string) + "',"
		}

	}
	buildsql = strings.TrimSuffix(buildsql, ",")
	endsql = strings.TrimSuffix(endsql, ",")
	buildsql = buildsql + ") VALUES (" + endsql + ");"

	id, _, err := DB.Execute(buildsql)
	if err != nil {
		return 0, err
	}

	return id, nil
}
