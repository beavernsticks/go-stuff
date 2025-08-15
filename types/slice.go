package bsgostuff_types

import (
	"fmt"
)

type Slice[T any] struct {
	values []T
	set    bool
	null   bool
}

func (s *Slice[T]) Set(value any) error {
	s.set = true
	s.null = false
	s.values = nil

	switch v := value.(type) {
	case []T:
		s.values = v
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}
