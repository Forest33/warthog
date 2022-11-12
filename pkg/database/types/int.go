// Package types provides basic operations with SQL types
package types

import (
	"database/sql"
)

// RefInt16ToSQL converts reference on int16 to sql.NullInt16
func RefInt16ToSQL(i *int16) sql.NullInt16 {
	if i == nil {
		return sql.NullInt16{}
	}
	return sql.NullInt16{Int16: *i, Valid: true}
}

// RefInt32ToSQL converts reference on int32 to sql.NullInt32
func RefInt32ToSQL(i *int32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: *i, Valid: true}
}

// RefInt64ToSQL converts reference on int64 to sql.NullInt64
func RefInt64ToSQL(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *i, Valid: true}
}

// SQLToRefInt16 converts sql.NullInt16 to reference on int16
func SQLToRefInt16(i sql.NullInt16) *int16 {
	if !i.Valid {
		return nil
	}
	return &i.Int16
}

// SQLToRefInt32 converts sql.NullInt32 to reference on int32
func SQLToRefInt32(i sql.NullInt32) *int32 {
	if !i.Valid {
		return nil
	}
	return &i.Int32
}

// SQLToRefInt64 converts sql.NullInt64 to reference on int64
func SQLToRefInt64(i sql.NullInt64) *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

// Int16ToSQL converts int16 to sql.NullInt16
func Int16ToSQL(i int16) sql.NullInt16 {
	return RefInt16ToSQL(&i)
}

// Int32ToSQL converts int32 to sql.NullInt32
func Int32ToSQL(i int32) sql.NullInt32 {
	return RefInt32ToSQL(&i)
}

// Int64ToSQL converts int16 to sql.NullInt64
func Int64ToSQL(i int64) sql.NullInt64 {
	return RefInt64ToSQL(&i)
}
