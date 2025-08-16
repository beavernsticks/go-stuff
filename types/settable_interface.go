package bsgostuff_types

// Settable defines the core interface for all nullable/optional types in the system.
// It provides methods for value manipulation, protocol buffers (gRPC) conversion,
// and PostgreSQL (pgx) integration with consistent nil-handling semantics.
//
// Implementations should maintain these invariants:
//   - IsNull() returns true only when explicitly set to null
//   - IsSet() distinguishes between unset and zero values
//   - GetValue() returns zero value when unset/null
//   - GetPtr() returns nil when unset/null
type Settable[T any, P any, D any] interface {
	// Set assigns a value from various representations.
	// Accepts nil, native type, proto wrapper, and pgx type.
	// Returns error for type mismatches.
	Set(any) error

	// GetValue returns the underlying value as a concrete type.
	// Returns zero value (e.g., "" for string) when unset/null.
	GetValue() T

	// GetPtr returns a pointer to the value.
	// Returns nil when unset/null to clearly indicate absence.
	GetPtr() *T

	// IsSet returns true if the value was explicitly set (even to null/zero).
	IsSet() bool

	// IsNull returns true if explicitly set to null.
	IsNull() bool

	// ToProto converts to protocol buffers representation.
	// Returns either concrete proto type or nil for unset/null.
	ToProto() P

	// ToPgx converts to pgx database type.
	// Returns pgx wrapper with Valid flag set appropriately.
	ToPgx() D
}
