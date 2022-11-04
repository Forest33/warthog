package types

import (
	"database/sql"
)

func RefStringToSQL(str *string) sql.NullString {
	if str == nil || *str == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: *str, Valid: true}
}

func RefEmptyStringToSQL(str *string) sql.NullString {
	if str == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *str, Valid: true}
}

func SQLToRefString(str sql.NullString) *string {
	if !str.Valid {
		return nil
	}
	return &str.String
}

func StringToSQL(str string) sql.NullString {
	return RefStringToSQL(&str)
}

func SQLToString(str sql.NullString) string {
	if !str.Valid {
		return ""
	}
	return str.String
}
