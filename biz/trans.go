package biz

import (
	"database/sql"
	"time"
)

// common type
type FieldType interface {
	int | int16 | int32 | int64 | uint | uint16 | uint32 | uint64 | float32 | float64 |
		bool |
		string | byte |
		time.Time
}

func PtrToType[T FieldType](from *T) T {
	var v T
	if from == nil {
		return v
	}

	return *from
}

func PtrToNullType[T FieldType](from *T) sql.Null[T] {
	if from == nil {
		return sql.Null[T]{
			Valid: false,
		}
	}
	return sql.Null[T]{
		Valid: true,
		V:     *from,
	}
}

func NullToPtrType[T FieldType](from sql.Null[T]) *T {
	if !from.Valid {
		return nil
	}
	return &from.V
}

func TypeToPtrType[T FieldType](from T) *T {
	return &from
}

// time
const defaultLayout = time.RFC3339

func TimePtrToTime(from *time.Time) time.Time {
	if from == nil {
		return time.UnixMilli(0) // default
	}

	return *from
}
func TimePtrToNullTime[T time.Time](from *time.Time) sql.Null[T] {
	if from == nil {
		return sql.Null[T]{
			Valid: false,
		} // default
	}

	return sql.Null[T]{
		Valid: true,
		V:     T(*from),
	}
}

func StrPtrToTime(from *string) time.Time {
	if from == nil {
		return time.UnixMilli(0) // default
	}

	t, err := time.Parse(defaultLayout, *from)
	if err != nil {
		return time.UnixMilli(0)
	}

	return t
}
func StrPtrToNullTime[T time.Time](from *string) sql.Null[T] {
	if from == nil {
		return sql.Null[T]{
			Valid: false,
		}
	}

	t, err := time.Parse(defaultLayout, *from)
	if err != nil {
		return sql.Null[T]{
			Valid: false,
		}
	}

	return sql.Null[T]{
		Valid: true,
		V:     T(t),
	}
}
func TimeToStrPtr(from time.Time) *string {
	s := from.Format(defaultLayout)
	return &s
}
func NullTimeToStrPtr(from sql.Null[time.Time]) *string {
	if !from.Valid {
		return nil
	}

	s := from.V.Format(defaultLayout)
	return &s
}
