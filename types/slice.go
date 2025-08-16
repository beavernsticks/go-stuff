package bsgostuff_types

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/proto"
)

// Slice implements a nullable generic slice type with three-state logic:
// - Set (with concrete elements)
// - Explicitly set to null
// - Unset (not initialized)
//
// Provides flexible slice handling with:
// - Protocol Buffers support (for repeated fields)
// - PostgreSQL array storage
// - JSON serialization/deserialization
//
// Implements Settable[[]T, []proto.Message, pgtype.Array[T]] where:
// - T is the element type
// - Proto representation is []proto.Message (for protobuf-compatible elements)
// - PostgreSQL storage uses pgtype.Array[T]
type Slice[T any] struct {
	items []T
	set   bool
	null  bool
}

// Compile-time interface check
var _ Settable[[]any, []proto.Message, pgtype.Array[any]] = (*Slice[any])(nil)

// NewSlice creates an initialized slice value.
// Returns concrete type for method chaining.
func NewSlice[T any](items []T) Slice[T] {
	return Slice[T]{items: items, set: true}
}

// Set assigns a value from supported types:
// - []T: direct slice value
// - pgtype.Array[T]: respects Valid flag
// - []proto.Message: for protobuf repeated fields
// - nil: explicit null
// - Other slice types: attempts element conversion
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
	case []proto.Message:
		items := make([]T, len(v))
		for i, msg := range v {
			if converted, ok := any(msg).(T); ok {
				items[i] = converted
			} else {
				return fmt.Errorf("element %d is not convertible to type %T", i, *new(T))
			}
		}
		s.items = items
	case nil:
		s.null = true
	default:
		// Attempt conversion via reflection
		val := reflect.ValueOf(value)
		if val.Kind() == reflect.Slice {
			items := make([]T, val.Len())
			for i := 0; i < val.Len(); i++ {
				elem := val.Index(i).Interface()
				if item, ok := elem.(T); ok {
					items[i] = item
				} else if msg, ok := elem.(proto.Message); ok {
					if converted, ok := any(msg).(T); ok {
						items[i] = converted
					} else {
						return fmt.Errorf("element %d is not convertible to type %T", i, *new(T))
					}
				} else {
					return fmt.Errorf("incompatible element type: %T", elem)
				}
			}
			s.items = items
		} else {
			return fmt.Errorf("unsupported type: %T", value)
		}
	}
	return nil
}

// GetValue returns the underlying slice.
// Returns nil when unset or null.
func (s Slice[T]) GetValue() []T {
	if !s.set || s.null {
		return nil
	}
	return s.items
}

// GetPtr returns a pointer to the slice.
// Returns nil when unset or null.
func (s Slice[T]) GetPtr() *[]T {
	if !s.set || s.null {
		return nil
	}
	return &s.items
}

// IsSet indicates whether the value was explicitly set.
// Returns false for uninitialized zero values.
func (s Slice[T]) IsSet() bool {
	return s.set
}

// IsNull indicates whether the value was explicitly set to null.
// Returns false for unset values.
func (s Slice[T]) IsNull() bool {
	return s.null
}

// ToProto converts to []proto.Message if elements are proto.Message.
// Returns nil when unset or null or for non-protobuf elements.
func (s Slice[T]) ToProto() []proto.Message {
	if !s.set || s.null {
		return nil
	}

	messages := make([]proto.Message, 0, len(s.items))
	for _, item := range s.items {
		if msg, ok := any(item).(proto.Message); ok {
			messages = append(messages, msg)
		} else {
			return nil
		}
	}
	return messages
}

// ToPgx converts to pgtype.Array[T] for PostgreSQL integration.
// Sets Valid flag according to null/unset state.
func (s Slice[T]) ToPgx() pgtype.Array[T] {
	return pgtype.Array[T]{
		Elements: s.items,
		Valid:    s.set && !s.null,
		Dims:     []pgtype.ArrayDimension{{Length: int32(len(s.items)), LowerBound: 1}},
	}
}

// MarshalJSON implements json.Marshaler interface.
func (s Slice[T]) MarshalJSON() ([]byte, error) {
	if !s.set || s.null {
		return []byte("null"), nil
	}
	return json.Marshal(s.items)
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (s *Slice[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		s.set = true
		s.null = true
		s.items = nil
		return nil
	}

	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	return s.Set(items)
}

// Append adds elements to the slice.
// Initializes the slice if not already set.
func (s *Slice[T]) Append(elements ...T) {
	s.set = true
	s.null = false
	s.items = append(s.items, elements...)
}

// Len returns the number of elements in the slice.
// Returns 0 when unset or null.
func (s Slice[T]) Len() int {
	if !s.set || s.null {
		return 0
	}
	return len(s.items)
}

// SliceMapper transforms a slice of type T to a slice of type R by applying a mapping function to each element.
//
// This is a generic implementation of the common map operation found in functional programming.
// It handles nil input safely and preserves the original slice length, applying the mapper
// function to every element in sequence.
//
// Example:
//
//	numbers := []int{1, 2, 3}
//	doubled := SliceMapper(numbers, func(n int) int {
//	    return n * 2
//	}) // returns []int{2, 4, 6}
//
//	nilSlice := []string(nil)
//	result := SliceMapper(nilSlice, strings.ToUpper) // returns nil
//
// Parameters:
//   - input: The input slice to transform (may be nil)
//   - mapper: Function that converts a T to an R
//
// Returns:
//   - A new slice with each element transformed by mapper
//   - nil if input is nil
func SliceMapper[T any, R any](input []T, mapper func(req T) R) []R {
	if input == nil {
		return nil
	}

	result := make([]R, len(input))
	for i, v := range input {
		result[i] = mapper(v)
	}

	return result
}

// SlicePtrMapper transforms a slice of pointers of type T to a slice of pointers of type R
// by applying a mapping function to each non-nil element.
//
// This function safely handles nil values in both the input slice and its elements.
// The mapper function is only called for non-nil elements, and nil elements in the input
// result in nil elements in the output at the same positions.
//
// Example:
//
//	users := []*User{user1, nil, user2}
//	names := SlicePtrMapper(users, func(u *User) *string {
//	    return &u.Name
//	}) // returns []*string{&user1.Name, nil, &user2.Name}
//
// Parameters:
//   - input: The input slice of pointers (may be nil or contain nil elements)
//   - mapper: Function that converts a *T to a *R
//
// Returns:
//   - A new slice with each non-nil element transformed by mapper
//   - nil if input is nil
//   - nil elements in positions where input had nil elements
func SlicePtrMapper[T any, R any](input []*T, mapper func(req *T) *R) []*R {
	if input == nil {
		return nil
	}

	result := make([]*R, len(input))
	for i, v := range input {
		if v != nil {
			result[i] = mapper(v)
		}
	}

	return result
}
