package bsgostuff_types

// ToPtr returns a pointer to the given value.
//
// This is useful for creating pointers to literals or values without needing
// to declare a separate variable. The function is generic and works with any type.
//
// Example:
//
//	strPtr := ToPtr("hello") // creates *string
//	intPtr := ToPtr(42)      // creates *int
//
// Parameters:
//   - val: The value to create a pointer for
//
// Returns:
//   - A pointer to the provided value
func ToPtr[T any](val T) *T {
	return &val
}

// Deref safely dereferences a pointer, returning either the pointed-to value
// or a default value if the pointer is nil.
//
// This provides nil-safety when working with pointers and avoids the need for
// explicit nil checks in calling code. The default value is always evaluated,
// even if not used (similar to how a normal function argument works).
//
// Example:
//
//	var ptr *int
//	value := Deref(ptr, 0) // returns 0
//	value := Deref(ToPtr(42), 0) // returns 42
//
// Parameters:
//   - ptr: The pointer to dereference (may be nil)
//   - defaultVal: The value to return if ptr is nil
//
// Returns:
//   - The dereferenced value or defaultVal if ptr is nil
func Deref[T any](ptr *T, defaultVal T) T {
	if ptr != nil {
		return *ptr
	}
	return defaultVal
}

// DerefZero safely dereferences a pointer, returning either the pointed-to value
// or the zero value for the type if the pointer is nil.
//
// This is similar to Deref but automatically uses the type's zero value
// instead of requiring an explicit default. Useful when you want to fall back
// to zero values rather than specific defaults.
//
// Example:
//
//	var ptr *string
//	value := DerefZero(ptr) // returns "" (zero value for string)
//
// Returns:
//   - The dereferenced value or the zero value for type T if ptr is nil
func DerefZero[T any](ptr *T) T {
	if ptr != nil {
		return *ptr
	}
	var zero T
	return zero
}

// DerefAny safely dereferences a pointer-like type, returning either the
// pointed-to value or a default if the pointer is nil.
//
// This works with any type that can be dereferenced with * operator, including
// custom pointer-like types (through the ~*E constraint). More flexible than
// Deref but requires slightly more complex type parameters.
//
// Example:
//
//	type MyPtr *int
//	var p MyPtr
//	value := DerefAny(p, 0) // returns 0
//
// Parameters:
//   - ptr: The pointer-like value to dereference (may be nil)
//   - defaultVal: The value to return if ptr is nil
//
// Returns:
//   - The dereferenced value or defaultVal if ptr is nil
func DerefAny[T ~*E, E any](ptr T, defaultVal E) E {
	if ptr != nil {
		return *ptr
	}
	return defaultVal
}

// DerefFunc conditionally processes a pointer value through a transformation
// function only if the pointer is non-nil.
//
// This provides a monadic-style operation for pointer chains, where the
// processor function is only called if the input pointer is non-nil. Useful
// for chaining operations on optional values.
//
// Example:
//
//	user := &User{Name: "Alice"}
//	namePtr := DerefFunc(user, func() *string { return &user.Name })
//
// Parameters:
//   - ptr: The pointer to check (may be nil)
//   - processor: Function to call if ptr is non-nil
//
// Returns:
//   - The result of processor() or nil if ptr was nil
func DerefFunc[T any, E any](ptr *T, processor func() *E) *E {
	if ptr == nil {
		return nil
	}
	return processor()
}
