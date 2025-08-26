package bsgostuff_types

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// String implements a nullable string with three-state logic:
// - Set (with concrete value)
// - Explicitly set to null
// - Unset (not initialized)
//
// Implements Settable[string, *wrapperspb.StringValue, pgtype.Text]
type String struct {
	value string
	set   bool
	null  bool
}

type Strings = Slice[string]

func NewStrings(items []string) Strings {
	return Strings(Slice[string]{items: items, set: true})
}

// Compile-time interface check
var _ Settable[string, *wrapperspb.StringValue, pgtype.Text] = (*String)(nil)

// NewString creates an initialized string value.
// Returns concrete type for method chaining.
func NewString(value string) String {
	return String{value: value, set: true}
}

// Set assigns a value from supported types:
// - string: direct value
// - *string: nil pointer treated as null
// - pgtype.Text: respects Valid flag
// - nil: explicit null
func (s *String) Set(value any) error {
	s.set = true
	s.null = false
	s.value = ""

	switch v := value.(type) {
	case string:
		s.value = v
	case *string:
		if v == nil {
			s.null = true
		} else {
			s.value = *v
		}
	case pgtype.Text:
		s.null = !v.Valid
		s.value = v.String
	case nil:
		s.null = true
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

// GetValue returns the underlying string value.
// Returns empty string when unset or null.
func (s String) GetValue() string {
	return s.value
}

// GetPtr returns a pointer to the string value.
// Returns nil when unset or null.
func (s String) GetPtr() *string {
	if !s.set || s.null {
		return nil
	}
	return &s.value
}

// IsSet indicates whether the value was explicitly set.
// Returns false for uninitialized zero values.
func (s String) IsSet() bool {
	return s.set
}

// IsNull indicates whether the value was explicitly set to null.
// Returns false for unset values.
func (s String) IsNull() bool {
	return s.null
}

// ToProto converts to protobuf StringValue wrapper.
// Returns nil when unset or null.
func (s String) ToProto() *wrapperspb.StringValue {
	if !s.set || s.null {
		return nil
	}
	return wrapperspb.String(s.value)
}

// ToPgx converts to pgtype.Text for PostgreSQL integration.
// Sets Valid flag according to null/unset state.
func (s String) ToPgx() pgtype.Text {
	return pgtype.Text{
		String: s.value,
		Valid:  s.set && !s.null,
	}
}
