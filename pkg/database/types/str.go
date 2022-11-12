// Package types provides basic operations with SQL types
package types

import (
	"database/sql"
)

// RefStringToSQL converts reference on string to sql.NullString
func RefStringToSQL(str *string) sql.NullString {
	if str == nil || *str == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: *str, Valid: true}
}

// RefEmptyStringToSQL converts reference on string to sql.NullString with check on empty
func RefEmptyStringToSQL(str *string) sql.NullString {
	if str == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *str, Valid: true}
}

// SQLToRefString converts sql.NullString to reference on string
func SQLToRefString(str sql.NullString) *string {
	if !str.Valid {
		return nil
	}
	return &str.String
}

// StringToSQL converts string to sql.NullString
func StringToSQL(str string) sql.NullString {
	return RefStringToSQL(&str)
}

// SQLToString converts sql.NullString to string
func SQLToString(str sql.NullString) string {
	if !str.Valid {
		return ""
	}
	return str.String
}
