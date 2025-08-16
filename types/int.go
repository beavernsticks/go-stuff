package bsgostuff_types

import (
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Int implements a nullable int64 with three-state logic:
// - Set (with concrete value)
// - Explicitly set to null
// - Unset (not initialized)
//
// Uses int64 as the underlying storage type for consistency and wide range.
// Provides conversion methods to int32 when needed.
//
// Implements Settable[int64, *wrapperspb.Int64Value, pgtype.Int8]
type Int struct {
	value int64
	set   bool
	null  bool
}

// Compile-time interface check
var _ Settable[int64, *wrapperspb.Int64Value, pgtype.Int8] = (*Int)(nil)

// NewInt creates an initialized int64 value.
// Returns concrete type for method chaining.
func NewInt(value int64) Int {
	return Int{value: value, set: true}
}

// Set assigns a value from supported types:
// - int, int32, int64: converted to int64
// - *int, *int32, *int64: nil pointer treated as null
// - string: parsed as int64 (empty string treated as zero)
// - pgtype.Int8: respects Valid flag
// - pgtype.Int4: respects Valid flag
// - nil: explicit null
func (i *Int) Set(value any) error {
	i.set = true
	i.null = false
	i.value = 0

	switch v := value.(type) {
	case int:
		i.value = int64(v)
	case int32:
		i.value = int64(v)
	case int64:
		i.value = v
	case *int:
		if v == nil {
			i.null = true
		} else {
			i.value = int64(*v)
		}
	case *int32:
		if v == nil {
			i.null = true
		} else {
			i.value = int64(*v)
		}
	case *int64:
		if v == nil {
			i.null = true
		} else {
			i.value = *v
		}
	case string:
		if v == "" {
			return nil // treat empty string as zero value
		}
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse string as int64: %w", err)
		}
		i.value = parsed
	case pgtype.Int8:
		i.null = !v.Valid
		i.value = v.Int64
	case pgtype.Int4:
		i.null = !v.Valid
		i.value = (int64)(v.Int32)
	case nil:
		i.null = true
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

// GetValue returns the underlying int64 value.
// Returns 0 when unset or null.
func (i Int) GetValue() int64 {
	return i.value
}

// GetPtr returns a pointer to the int64 value.
// Returns nil when unset or null.
func (i Int) GetPtr() *int64 {
	if !i.set || i.null {
		return nil
	}
	return &i.value
}

// IsSet indicates whether the value was explicitly set.
// Returns false for uninitialized zero values.
func (i Int) IsSet() bool {
	return i.set
}

// IsNull indicates whether the value was explicitly set to null.
// Returns false for unset values.
func (i Int) IsNull() bool {
	return i.null
}

// ToProto converts to protobuf Int64Value wrapper.
// Returns nil when unset or null.
func (i Int) ToProto() *wrapperspb.Int64Value {
	if !i.set || i.null {
		return nil
	}
	return wrapperspb.Int64(i.value)
}

// ToPgx converts to pgtype.Int8 for PostgreSQL integration.
// Sets Valid flag according to null/unset state.
func (i Int) ToPgx() pgtype.Int8 {
	return pgtype.Int8{
		Int64: i.value,
		Valid: i.set && !i.null,
	}
}
