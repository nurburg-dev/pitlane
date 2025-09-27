package db

import (
	"errors"
	"reflect"

	"github.com/jackc/pgx/v5"
)

// RowMapper provides generic mapping functionality for pgx results
type RowMapper struct{}

func NewRowMapper() *RowMapper {
	return &RowMapper{}
}

// ScanRow scans a single row into a destination struct using reflection
func (rm *RowMapper) ScanRow(row pgx.Row, dest interface{}) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr || destValue.Elem().Kind() != reflect.Struct {
		return errors.New("dest must be a pointer to struct")
	}

	structValue := destValue.Elem()
	structType := structValue.Type()

	// Create slice of pointers to struct fields
	fieldPtrs := make([]interface{}, structType.NumField())
	for i := range structType.NumField() {
		fieldPtrs[i] = structValue.Field(i).Addr().Interface()
	}

	return row.Scan(fieldPtrs...)
}

// ScanRows scans multiple rows into a slice of structs using reflection
func (rm *RowMapper) ScanRows(rows pgx.Rows, dest interface{}) error {
	defer rows.Close()

	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr || destValue.Elem().Kind() != reflect.Slice {
		return errors.New("dest must be a pointer to slice")
	}

	sliceValue := destValue.Elem()
	elementType := sliceValue.Type().Elem()

	if elementType.Kind() != reflect.Struct {
		return errors.New("slice elements must be structs")
	}

	for rows.Next() {
		// Create new instance of struct
		newElem := reflect.New(elementType).Elem()

		// Create slice of pointers to struct fields
		fieldPtrs := make([]any, elementType.NumField())
		for i := range elementType.NumField() {
			fieldPtrs[i] = newElem.Field(i).Addr().Interface()
		}

		if err := rows.Scan(fieldPtrs...); err != nil {
			return err
		}

		// Append to slice
		sliceValue.Set(reflect.Append(sliceValue, newElem))
	}

	return rows.Err()
}
