package bsgostuff_proto

import (
	bsgostuff_domain "github.com/beavernsticks/go-stuff/domain"
	bsgostuff_types "github.com/beavernsticks/go-stuff/types"
)

func ToSortProto(value *bsgostuff_domain.Sort) *Sort {
	_sort := bsgostuff_types.DerefZero(value)

	return &Sort{
		Field:   _sort.Field.GetPtr(),
		Reverse: _sort.IsReversed.GetPtr(),
	}
}

func FromSortProto(value *Sort) bsgostuff_domain.Sort {
	if value == nil {
		return bsgostuff_domain.Sort{}
	}

	_field := bsgostuff_types.String{}
	if value.Field != nil {
		_field.Set(*value.Field)
	}

	_isReversed := bsgostuff_types.Bool{}
	if value.Reverse != nil {
		_isReversed.Set(*value.Reverse)
	}

	return bsgostuff_domain.Sort{
		Field:      _field,
		IsReversed: _isReversed,
	}
}

func (value *Sort) FromProto() bsgostuff_domain.Sort {
	return FromSortProto(value)
}
