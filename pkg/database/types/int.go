package types

import (
	"database/sql"
)

func RefInt16ToSQL(i *int16) sql.NullInt16 {
	if i == nil {
		return sql.NullInt16{}
	}
	return sql.NullInt16{Int16: *i, Valid: true}
}

func RefInt32ToSQL(i *int32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: *i, Valid: true}
}

func RefInt64ToSQL(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *i, Valid: true}
}

func SQLToRefInt16(i sql.NullInt16) *int16 {
	if !i.Valid {
		return nil
	}
	return &i.Int16
}

func SQLToRefInt32(i sql.NullInt32) *int32 {
	if !i.Valid {
		return nil
	}
	return &i.Int32
}

func SQLToRefInt64(i sql.NullInt64) *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

func Int16ToSQL(i int16) sql.NullInt16 {
	return RefInt16ToSQL(&i)
}

func Int32ToSQL(i int32) sql.NullInt32 {
	return RefInt32ToSQL(&i)
}

func Int64ToSQL(i int64) sql.NullInt64 {
	return RefInt64ToSQL(&i)
}
