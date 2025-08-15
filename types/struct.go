package bsgostuff_types

import (
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type Struct[T any] struct {
	value *T
	set   bool
	null  bool
}

func NewStruct[T any](v T) Struct[T] {
	return Struct[T]{value: &v, set: true}
}

func (s *Struct[T]) Set(value any) error {
	s.set = true
	s.null = false
	s.value = nil

	switch v := value.(type) {
	case T:
		s.value = &v
	case *T:
		s.value = v
		s.null = (v == nil)
	case pgtype.Text:
		if !v.Valid {
			s.null = true
		} else {
			var parsed T
			if err := json.Unmarshal([]byte(v.String), &parsed); err != nil {
				return err
			}
			s.value = &parsed
		}
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

func (s Struct[T]) ToPgx() pgtype.Text {
	if !s.set || s.null {
		return pgtype.Text{Valid: false}
	}
	data, _ := json.Marshal(s.value)
	return pgtype.Text{String: string(data), Valid: true}
}
