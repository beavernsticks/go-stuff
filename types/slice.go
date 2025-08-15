package bsgostuff_types

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v5/pgtype"
)

// Slice обобщённый тип для массивов с поддержкой nil/PGX
type Slice[T any] struct {
	items []T
	set   bool
	null  bool
}

// NewSlice конструктор
func NewSlice[T any](items []T) Slice[T] {
	return Slice[T]{items: items, set: true}
}

// Set принимает: []T, pgtype.Array[T], nil
func (s *Slice[T]) Set(value any) error {
	s.set = true
	s.null = false
	s.items = nil

	switch v := value.(type) {
	case []T:
		s.items = v
	case pgtype.Array[T]:
		if !v.Valid {
			s.null = true
		} else {
			s.items = v.Elements
		}
	case nil:
		s.null = true
	default:
		// Попытка конвертации через reflection
		val := reflect.ValueOf(value)
		if val.Kind() == reflect.Slice {
			items := make([]T, val.Len())
			for i := 0; i < val.Len(); i++ {
				if item, ok := val.Index(i).Interface().(T); ok {
					items[i] = item
				} else {
					return fmt.Errorf("несовместимый тип элемента: %T", val.Index(i).Interface())
				}
			}
			s.items = items
		} else {
			return fmt.Errorf("неподдерживаемый тип: %T", value)
		}
	}
	return nil
}

// Get возвращает []T (nil если null/unset)
func (s Slice[T]) Get() []T {
	if !s.set || s.null {
		return nil
	}
	return s.items
}

// ToPgx конвертация в pgtype.Array[T]
func (s Slice[T]) ToPgx() pgtype.Array[T] {
	return pgtype.Array[T]{
		Elements: s.items,
		Valid:    s.set && !s.null,
	}
}

// -- Стандартные методы --
func (s Slice[T]) IsSet() bool  { return s.set }
func (s Slice[T]) IsNull() bool { return s.null }

// MarshalJSON для JSON
func (s Slice[T]) MarshalJSON() ([]byte, error) {
	if !s.set || s.null {
		return []byte("null"), nil
	}
	return json.Marshal(s.items)
}
