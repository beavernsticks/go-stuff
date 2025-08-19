package bsgostuff_types

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Timestamp implements a nullable time.Time with three-state logic:
// - Set (with concrete value)
// - Explicitly set to null
// - Unset (not initialized)
//
// Provides comprehensive time handling with:
// - RFC3339 string parsing
// - Unix epoch conversions
// - Protobuf Timestamp integration
// - PostgreSQL timestamp storage
//
// Implements Settable[time.Time, *timestamppb.Timestamp, pgtype.Timestamp]
type Timestamp struct {
	value time.Time
	set   bool
	null  bool
}

type Timestamps = Slice[Timestamp]

func NewTimestamps(items []Timestamp) Timestamps {
	return Timestamps(Slice[Timestamp]{items: items, set: true})
}

// Compile-time interface check
var _ Settable[time.Time, *timestamppb.Timestamp, pgtype.Timestamp] = (*Timestamp)(nil)

// NewTimestamp creates an initialized timestamp value.
// Returns concrete type for method chaining.
func NewTimestamp(value time.Time) Timestamp {
	return Timestamp{value: value, set: true}
}

// Set assigns a value from supported types:
// - time.Time: direct value
// - *time.Time: nil pointer treated as null
// - string: parsed as RFC3339 or Unix timestamp
// - int64: treated as Unix timestamp (seconds since epoch)
// - pgtype.Timestamp: respects Valid flag
// - timestamppb.Timestamp: protobuf timestamp
// - nil: explicit null
func (t *Timestamp) Set(value any) error {
	t.set = true
	t.null = false
	t.value = time.Time{}

	switch v := value.(type) {
	case time.Time:
		t.value = v
	case *time.Time:
		if v == nil {
			t.null = true
		} else {
			t.value = *v
		}
	case string:
		if v == "" {
			t.null = true
			return nil
		}

		// Try RFC3339 first
		parsed, err := time.Parse(time.RFC3339, v)
		if err == nil {
			t.value = parsed
			return nil
		}

		// Fall back to Unix timestamp
		unix, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse timestamp string (neither RFC3339 nor Unix timestamp): %w", err)
		}
		t.value = time.Unix(unix, 0)
	case int64:
		t.value = time.Unix(v, 0)
	case pgtype.Timestamp:
		t.null = !v.Valid
		t.value = v.Time
	case *timestamppb.Timestamp:
		if v == nil {
			t.null = true
		} else {
			t.value = v.AsTime()
		}
	case timestamppb.Timestamp:
		t.value = v.AsTime()
	case nil:
		t.null = true
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

// GetValue returns the underlying time.Time value.
// Returns zero time when unset or null.
func (t Timestamp) GetValue() time.Time {
	return t.value
}

// GetPtr returns a pointer to the time.Time value.
// Returns nil when unset or null.
func (t Timestamp) GetPtr() *time.Time {
	if !t.set || t.null {
		return nil
	}
	return &t.value
}

// IsSet indicates whether the value was explicitly set.
// Returns false for uninitialized zero values.
func (t Timestamp) IsSet() bool {
	return t.set
}

// IsNull indicates whether the value was explicitly set to null.
// Returns false for unset values.
func (t Timestamp) IsNull() bool {
	return t.null
}

// ToProto converts to protobuf Timestamp.
// Returns nil when unset or null.
func (t Timestamp) ToProto() *timestamppb.Timestamp {
	if !t.set || t.null {
		return nil
	}
	return timestamppb.New(t.value)
}

// ToPgx converts to pgtype.Timestamp for PostgreSQL integration.
// Sets Valid flag according to null/unset state.
func (t Timestamp) ToPgx() pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  t.value,
		Valid: t.set && !t.null,
	}
}

// Unix returns the Unix timestamp (seconds since January 1, 1970 UTC).
// Returns 0 when unset or null.
func (t Timestamp) Unix() int64 {
	if !t.set || t.null {
		return 0
	}
	return t.value.Unix()
}

// UnixNano returns the Unix timestamp in nanoseconds.
// Returns 0 when unset or null.
func (t Timestamp) UnixNano() int64 {
	if !t.set || t.null {
		return 0
	}
	return t.value.UnixNano()
}

// Format returns a textual representation of the timestamp.
// Returns empty string when unset or null.
func (t Timestamp) Format(layout string) string {
	if !t.set || t.null {
		return ""
	}
	return t.value.Format(layout)
}

// RFC3339 returns the timestamp in RFC3339 format.
// Returns empty string when unset or null.
func (t Timestamp) RFC3339() string {
	return t.Format(time.RFC3339)
}

// NewTimestampFromUnix creates a new Timestamp from Unix timestamp (seconds since epoch).
func NewTimestampFromUnix(sec int64) Timestamp {
	return NewTimestamp(time.Unix(sec, 0))
}

// NewTimestampFromUnixNano creates a new Timestamp from Unix nanosecond timestamp.
func NewTimestampFromUnixNano(nsec int64) Timestamp {
	return NewTimestamp(time.Unix(0, nsec))
}

// NewCurrentTimestamp creates a new Timestamp with current time.
func NewCurrentTimestamp() Timestamp {
	return NewTimestamp(time.Now())
}
