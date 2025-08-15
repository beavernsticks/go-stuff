package bsgostuff_types

import (
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v5/pgtype"
)

type Int struct {
	value int64
	set   bool
	null  bool
}

func (i *Int) Set(value any) error {
	i.set = true
	i.null = false

	switch v := value.(type) {
	case int, int32, int64:
		i.value = reflect.ValueOf(v).Int()
	case *int64:
		if v == nil {
			i.null = true
		} else {
			i.value = *v
		}
	case pgtype.Int4:
		if v.Valid {
			i.value = int64(v.Int32)
		} else {
			i.null = true
		}
	case pgtype.Int8:
		if v.Valid {
			i.value = v.Int64
		} else {
			i.null = true
		}
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

func (i Int) Get() *int64 {
	if !i.set || i.null {
		return nil
	}
	return &i.value
}
