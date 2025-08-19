package bsgostuff_types

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/structpb"
)

// JSON implements a nullable JSON container with three-state logic:
// - Set (with concrete JSON value)
// - Explicitly set to null
// - Unset (not initialized)
//
// Handles arbitrary JSON data with:
// - Efficient storage as raw bytes
// - Lazy unmarshaling
// - Support for both structured and unstructured data
// - Protobuf Struct/Value integration
// - PostgreSQL jsonb storage
//
// Implements Settable[json.RawMessage, *structpb.Struct, pgtype.Text]
type JSON struct {
	data  json.RawMessage
	set   bool
	null  bool
	valid bool // tracks whether data is valid JSON
}

type JSONs = Slice[JSON]

func NewJSONs(items []JSON) JSONs {
	return JSONs(Slice[JSON]{items: items, set: true})
}

// Compile-time interface check
var _ Settable[json.RawMessage, *structpb.Struct, pgtype.Text] = (*JSON)(nil)

// NewJSON creates an initialized JSON value from raw JSON bytes.
// Returns concrete type for method chaining.
func NewJSON(data []byte) JSON {
	j := JSON{}
	j.Set(data) // handles validation
	return j
}

// NewJSONFromValue creates a JSON value by marshaling a Go value.
// Returns error if marshaling fails.
func NewJSONFromValue(v any) (JSON, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return JSON{}, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return NewJSON(data), nil
}

// MustNewJSONFromValue creates a JSON value by marshaling a Go value.
// Panics if marshaling fails.
func MustNewJSONFromValue(v any) JSON {
	j, err := NewJSONFromValue(v)
	if err != nil {
		panic(err)
	}
	return j
}

// Set assigns a value from supported types:
// - []byte: raw JSON (must be valid JSON)
// - string: JSON string (must be valid JSON)
// - json.RawMessage: raw JSON
// - any Go value: will be marshaled to JSON
// - pgtype.Text: respects Valid flag
// - *structpb.Struct: protobuf Struct
// - *structpb.Value: protobuf Value
// - nil: explicit null
func (j *JSON) Set(value any) error {
	j.set = true
	j.null = false
	j.valid = false
	j.data = nil

	switch v := value.(type) {
	case []byte:
		if !json.Valid(v) {
			return fmt.Errorf("invalid JSON data")
		}
		j.data = v
		j.valid = true
	case string:
		if !json.Valid([]byte(v)) {
			return fmt.Errorf("invalid JSON string")
		}
		j.data = []byte(v)
		j.valid = true
	case json.RawMessage:
		if !json.Valid(v) {
			return fmt.Errorf("invalid json.RawMessage")
		}
		j.data = v
		j.valid = true
	case pgtype.Text:
		j.null = !v.Valid
		if v.Valid {
			if !json.Valid([]byte(v.String)) {
				return fmt.Errorf("invalid JSON in pgtype.Text")
			}
			j.data = []byte(v.String)
			j.valid = true
		}
	case *structpb.Struct:
		if v == nil {
			j.null = true
		} else {
			data, err := v.MarshalJSON()
			if err != nil {
				return fmt.Errorf("failed to marshal protobuf Struct: %w", err)
			}
			j.data = data
			j.valid = true
		}
	case *structpb.Value:
		if v == nil {
			j.null = true
		} else {
			data, err := v.MarshalJSON()
			if err != nil {
				return fmt.Errorf("failed to marshal protobuf Value: %w", err)
			}
			j.data = data
			j.valid = true
		}
	case nil:
		j.null = true
	default:
		// Attempt to marshal arbitrary Go value
		data, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal value to JSON: %w", err)
		}
		j.data = data
		j.valid = true
	}
	return nil
}

// GetValue returns the underlying JSON as raw bytes.
// Returns nil when unset or null.
func (j JSON) GetValue() json.RawMessage {
	if !j.set || j.null {
		return nil
	}
	return j.data
}

// GetPtr returns a pointer to the raw JSON data.
// Returns nil when unset or null.
func (j JSON) GetPtr() *json.RawMessage {
	if !j.set || j.null {
		return nil
	}
	return &j.data
}

// IsSet indicates whether the value was explicitly set.
// Returns false for uninitialized zero values.
func (j JSON) IsSet() bool {
	return j.set
}

// IsNull indicates whether the value was explicitly set to null.
// Returns false for unset values.
func (j JSON) IsNull() bool {
	return j.null
}

// IsValid indicates whether the contained data is valid JSON.
// Returns false when unset, null, or contains invalid JSON.
func (j JSON) IsValid() bool {
	return j.valid
}

// ToProto converts to protobuf Struct or Value.
// Returns nil when unset or null.
func (j JSON) ToProto() *structpb.Struct {
	if !j.set || j.null || !j.valid {
		return nil
	}

	var st structpb.Struct
	if err := st.UnmarshalJSON(j.data); err != nil {
		return nil
	}
	return &st
}

// ToPgx converts to pgtype.Text for PostgreSQL integration.
// Sets Valid flag according to null/unset state.
func (j JSON) ToPgx() pgtype.Text {
	if !j.set || j.null {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{
		String: string(j.data),
		Valid:  true,
	}
}

// Unmarshal parses the JSON-encoded data into the value pointed to by v.
// Returns error if JSON is invalid or unmarshaling fails.
func (j JSON) Unmarshal(v any) error {
	if !j.set || j.null {
		return fmt.Errorf("cannot unmarshal null or unset JSON")
	}
	if !j.valid {
		return fmt.Errorf("contains invalid JSON data")
	}
	return json.Unmarshal(j.data, v)
}

// MarshalJSON implements json.Marshaler interface.
func (j JSON) MarshalJSON() ([]byte, error) {
	if !j.set || j.null {
		return []byte("null"), nil
	}
	if !j.valid {
		return nil, fmt.Errorf("cannot marshal invalid JSON")
	}
	return j.data, nil
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (j *JSON) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		j.set = true
		j.null = true
		j.valid = true
		j.data = nil
		return nil
	}

	if !json.Valid(data) {
		return fmt.Errorf("invalid JSON data")
	}

	j.set = true
	j.null = false
	j.valid = true
	j.data = data
	return nil
}
