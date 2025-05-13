package gsdb

import (
	_ "errors"
	"fmt"
	"reflect"
	"strings"
)

// Insert generates an SQL query based on the db column tags provided in the structure of the argument
func (db *Database) Insert(dbStructure any) (string, error) {
	t := reflect.TypeOf(dbStructure)
	table, buildSql, err := generateBuildSql(dbStructure, t)
	if err != nil {
		return "", err
	}
	if table == "" {
		return "", fmt.Errorf("no table found in structure")
	}
	if buildSql == "" {
		return "", fmt.Errorf("no non-primary key and non-omitted fields found in structure")
	}
	valueSql, err := generateValuesSql(dbStructure, t)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("INSERT INTO %s(%s) VALUES %s;", table, buildSql, valueSql), nil
}

// InsertMany generates an SQL query based on the db column tags provided in the structure of the elements in the argument
func InsertMany[T any](dbStructures []T) (string, error) {
	if len(dbStructures) == 0 {
		return "", nil
	}
	t := reflect.TypeOf(dbStructures[0])
	table, buildSql, err := generateBuildSql(dbStructures[0], t)
	if err != nil {
		return "", err
	}
	if table == "" {
		return "", fmt.Errorf("no table found in structure")
	}
	if buildSql == "" {
		return "", fmt.Errorf("no non-primary key and non-omitted fields found in structure")
	}
	var valuesSql strings.Builder
	entriesLength := len(dbStructures)
	for i, dbStructure := range dbStructures {
		valueSql, err := generateValuesSql(dbStructure, t)
		if err != nil {
			return "", err
		}
		valuesSql.WriteString(valueSql)
		if i < entriesLength-1 {
			valuesSql.WriteString("\n")
		}
	}
	return fmt.Sprintf("INSERT INTO %s(%s) VALUES %s;", table, buildSql, valuesSql.String()), nil
}
