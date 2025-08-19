package bsgostuff_types

import (
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Bool implements a nullable boolean with three-state logic:
// - Set (true/false)
// - Explicitly set to null
// - Unset (not initialized)
//
// Implements Settable[bool, *wrapperspb.BoolValue, pgtype.Bool] interface
// for full type safety across all conversions.
type Bool struct {
	value bool
	set   bool
	null  bool
}

type Bools = Slice[Bool]

func NewBools(items []Bool) Bools {
	return Bools(Slice[Bool]{items: items, set: true})
}

// Compile-time interface check
var _ Settable[bool, *wrapperspb.BoolValue, pgtype.Bool] = (*Bool)(nil)

// NewBool creates an initialized boolean value.
// Returns concrete type for better method chaining.
func NewBool(value bool) Bool {
	return Bool{value: value, set: true}
}

// Set assigns a value from supported types:
// - bool: direct value
// - *bool: nil pointer treated as null
// - string: "true"/"1" or "false"/"0"
// - *string: nil treated as null
// - pgtype.Bool: respects Valid flag
// - nil: explicit null
func (b *Bool) Set(value any) error {
	b.set = true
	b.null = false
	b.value = false

	switch v := value.(type) {
	case bool:
		b.value = v
	case *bool:
		if v == nil {
			b.null = true
		} else {
			b.value = *v
		}
	case string:
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			return fmt.Errorf("invalid bool string: %w", err)
		}
		b.value = parsed
	case *string:
		if v == nil {
			b.null = true
		} else {
			parsed, err := strconv.ParseBool(*v)
			if err != nil {
				return fmt.Errorf("invalid bool string: %w", err)
			}
			b.value = parsed
		}
	case pgtype.Bool:
		b.null = !v.Valid
		b.value = v.Bool
	case nil:
		b.null = true
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

// GetValue returns the underlying boolean value.
// Returns false when unset or null (zero value for bool).
func (b Bool) GetValue() bool {
	return b.value
}

// GetPtr returns a pointer to the boolean value.
// Returns nil when unset or null.
func (b Bool) GetPtr() *bool {
	if !b.set || b.null {
		return nil
	}
	return &b.value
}

// IsSet indicates whether the value was explicitly set.
// Returns false for uninitialized zero values.
func (b Bool) IsSet() bool {
	return b.set
}

// IsNull indicates whether the value was explicitly set to null.
// Returns false for unset values.
func (b Bool) IsNull() bool {
	return b.null
}

// ToProto converts to protobuf BoolValue wrapper.
// Returns nil when unset or null.
func (b Bool) ToProto() *wrapperspb.BoolValue {
	if !b.set || b.null {
		return nil
	}
	return wrapperspb.Bool(b.value)
}

// ToPgx converts to pgtype.Bool for PostgreSQL integration.
// Sets Valid flag according to null/unset state.
func (b Bool) ToPgx() pgtype.Bool {
	return pgtype.Bool{
		Bool:  b.value,
		Valid: b.set && !b.null,
	}
}
