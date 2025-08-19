package bsgostuff_types

import (
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Struct implements a nullable generic struct type with three-state logic:
// - Set (with concrete value)
// - Explicitly set to null
// - Unset (not initialized)
//
// Provides JSON serialization/deserialization for structs with:
// - Protocol Buffers support (for proto.Message types)
// - Standard JSON marshaling
// - PostgreSQL text/jsonb storage
//
// Implements Settable[T, proto.Message, pgtype.Text] where:
// - T is the struct type
// - Proto representation is the same as native (for proto.Message)
// - PostgreSQL storage uses pgtype.Text (stored as jsonb in database)
type Struct[T any] struct {
	value *T
	set   bool
	null  bool
}

type Structs[T any] = Slice[Struct[T]]

func NewStructs[T any](items []Struct[T]) Structs[T] {
	return Structs[T](Slice[Struct[T]]{items: items, set: true})
}

// Compile-time interface check
var _ Settable[any, proto.Message, pgtype.Text] = (*Struct[any])(nil)

// NewStruct creates an initialized struct value.
// Returns concrete type for method chaining.
func NewStruct[T any](value T) Struct[T] {
	return Struct[T]{value: &value, set: true}
}

// Set assigns a value from supported types:
// - Native struct type (T): direct value
// - *T: nil pointer treated as null
// - string: parsed as JSON
// - []byte: parsed as JSON
// - pgtype.Text: respects Valid flag
// - proto.Message: for protobuf structs
// - nil: explicit null
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
	case string:
		if v == "" {
			s.null = true
			return nil
		}
		return s.unmarshalJSON([]byte(v))
	case []byte:
		if len(v) == 0 {
			s.null = true
			return nil
		}
		return s.unmarshalJSON(v)
	case pgtype.Text:
		if !v.Valid {
			s.null = true
		} else {
			return s.unmarshalJSON([]byte(v.String))
		}
	case proto.Message:
		// Special handling for protobuf messages
		data, err := protojson.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal protobuf: %w", err)
		}
		return s.unmarshalJSON(data)
	case nil:
		s.null = true
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

// unmarshalJSON handles JSON deserialization with special protobuf support
func (s *Struct[T]) unmarshalJSON(data []byte) error {
	var zero T

	// Check if the type implements proto.Message
	if _, isProto := any(zero).(proto.Message); isProto {
		msg := any(zero).(proto.Message)
		if err := protojson.Unmarshal(data, msg); err != nil {
			return fmt.Errorf("failed to unmarshal protobuf: %w", err)
		}
		s.value = any(msg).(*T)
		return nil
	}

	// Standard JSON unmarshaling
	var parsed T
	if err := json.Unmarshal(data, &parsed); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	s.value = &parsed
	return nil
}

// GetValue returns the underlying struct value.
// Returns zero value when unset or null.
func (s Struct[T]) GetValue() T {
	if !s.set || s.null {
		var zero T
		return zero
	}
	return *s.value
}

// GetPtr returns a pointer to the struct value.
// Returns nil when unset or null.
func (s Struct[T]) GetPtr() *T {
	if !s.set || s.null {
		return nil
	}
	return s.value
}

// IsSet indicates whether the value was explicitly set.
// Returns false for uninitialized zero values.
func (s Struct[T]) IsSet() bool {
	return s.set
}

// IsNull indicates whether the value was explicitly set to null.
// Returns false for unset values.
func (s Struct[T]) IsNull() bool {
	return s.null
}

// ToProto converts to protobuf Message if the struct is a proto.Message.
// Returns nil when unset or null or for non-protobuf structs.
func (s Struct[T]) ToProto() proto.Message {
	if !s.set || s.null {
		return nil
	}

	if msg, ok := any(*s.value).(proto.Message); ok {
		return msg
	}
	return nil
}

// ToPgx converts to pgtype.Text for PostgreSQL integration.
// Uses JSON serialization and sets Valid flag according to null/unset state.
func (s Struct[T]) ToPgx() pgtype.Text {
	if !s.set || s.null {
		return pgtype.Text{Valid: false}
	}

	var data []byte
	var err error

	// Use protojson for protobuf messages
	if msg, ok := any(*s.value).(proto.Message); ok {
		data, err = protojson.Marshal(msg)
	} else {
		data, err = json.Marshal(s.value)
	}

	if err != nil {
		return pgtype.Text{Valid: false}
	}

	return pgtype.Text{
		String: string(data),
		Valid:  true,
	}
}

// MarshalJSON implements json.Marshaler interface
func (s Struct[T]) MarshalJSON() ([]byte, error) {
	if !s.set || s.null {
		return []byte("null"), nil
	}
	return json.Marshal(s.value)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (s *Struct[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		s.set = true
		s.null = true
		s.value = nil
		return nil
	}
	return s.Set(data)
}
