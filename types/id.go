package bsgostuff_types

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// ID implements a nullable UUID identifier with three-state logic:
// - Set (with concrete UUID value)
// - Explicitly set to null
// - Unset (not initialized)
//
// Implements Settable[uuid.UUID, *wrapperspb.StringValue, pgtype.UUID]
type ID struct {
	value uuid.UUID
	set   bool
	null  bool
}

// Compile-time interface check
var _ Settable[uuid.UUID, *wrapperspb.StringValue, pgtype.UUID] = (*ID)(nil)

// NewID creates a new initialized UUID value.
// Generates random UUID using google/uuid package.
func NewID() ID {
	return ID{
		value: uuid.New(),
		set:   true,
	}
}

// NewIDFromString creates an ID from UUID string.
// Returns error if string is not valid UUID format.
func NewIDFromString(s string) (ID, error) {
	id := ID{}
	err := id.Set(s)
	return id, err
}

// Set assigns a value from supported types:
// - string: must be valid UUID format
// - *string: nil pointer treated as null
// - uuid.UUID: direct UUID value
// - pgtype.UUID: respects Valid flag
// - nil: explicit null
func (i *ID) Set(value any) error {
	i.set = true
	i.null = false
	i.value = uuid.Nil // Zero value

	switch v := value.(type) {
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("invalid UUID format: %w", err)
		}
		i.value = parsed
	case *string:
		if v == nil {
			i.null = true
		} else {
			parsed, err := uuid.Parse(*v)
			if err != nil {
				return fmt.Errorf("invalid UUID format: %w", err)
			}
			i.value = parsed
		}
	case uuid.UUID:
		i.value = v
	case pgtype.UUID:
		i.null = !v.Valid
		i.value = v.Bytes
	case nil:
		i.null = true
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

// GetValue returns the underlying UUID value.
// Returns uuid.Nil when unset or null.
func (i ID) GetValue() uuid.UUID {
	return i.value
}

// GetPtr returns a pointer to the UUID value.
// Returns nil when unset or null.
func (i ID) GetPtr() *uuid.UUID {
	if !i.set || i.null {
		return nil
	}
	return &i.value
}

// IsSet indicates whether the value was explicitly set.
// Returns false for uninitialized zero values.
func (i ID) IsSet() bool {
	return i.set
}

// IsNull indicates whether the value was explicitly set to null.
// Returns false for unset values.
func (i ID) IsNull() bool {
	return i.null
}

// ToProto converts to protobuf StringValue wrapper.
// Returns nil when unset or null.
// UUID is transmitted as canonical string representation.
func (i ID) ToProto() *wrapperspb.StringValue {
	if !i.set || i.null {
		return nil
	}
	return wrapperspb.String(i.value.String())
}

// ToPgx converts to pgtype.UUID for PostgreSQL integration.
// Sets Valid flag according to null/unset state.
func (i ID) ToPgx() pgtype.UUID {
	return pgtype.UUID{
		Bytes: i.value,
		Valid: i.set && !i.null,
	}
}

// String returns canonical UUID string representation.
// Returns empty string when unset or null.
func (i ID) String() string {
	if !i.set || i.null {
		return ""
	}
	return i.value.String()
}
