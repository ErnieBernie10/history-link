package db

import (
	"database/sql"
	"time"

	"github.com/huandu/go-sqlbuilder"
)

func ScanStructRow[T any](row *sql.Row, o *T, structure *sqlbuilder.Struct) error {
	if err := row.Scan(structure.Addr(&o)...); err != nil {
		return err
	}
	return nil
}

func ScanStructRows[T any](rows *sql.Rows, l *[]T, structure *sqlbuilder.Struct) error {
	var objects []T
	for rows.Next() {
		var object T
		if err := rows.Scan(structure.Addr(&object)...); err != nil {
			return err
		}
		objects = append(objects, object)
	}
	*l = objects
	return nil
}

func NewNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NewNullInt16(i int16) sql.NullInt16 {
	return sql.NullInt16{
		Int16: i,
		Valid: true,
	}
}

func NewNullInt32(i int32) sql.NullInt32 {
	return sql.NullInt32{
		Int32: i,
		Valid: true,
	}
}

func NewNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: !t.IsZero(),
	}
}
