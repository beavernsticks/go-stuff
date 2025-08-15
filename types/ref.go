package bsgostuff_types

func ToPtr[T any](val T) *T {
	return &val
}

// Deref возвращает разыменованное значение или defaultVal, если ptr == nil.
func Deref[T any](ptr *T, defaultVal T) T {
	if ptr != nil {
		return *ptr
	}
	return defaultVal
}

func DerefZero[T any](ptr *T) T {
	if ptr != nil {
		return *ptr
	}
	var zero T
	return zero
}

func DerefAny[T ~*E, E any](ptr T, defaultVal E) E {
	if ptr != nil {
		return *ptr
	}
	return defaultVal
}

func DerefFunc[T any, E any](value *T, processor func() *E) *E {
	if value == nil {
		return nil
	}

	return processor()
}
