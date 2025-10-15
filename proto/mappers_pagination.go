package bsgostuff_proto

import (
	bsgostuff_domain "github.com/beavernsticks/go-stuff/domain"
	bsgostuff_types "github.com/beavernsticks/go-stuff/types"
)

func ToPaginationProto(value *bsgostuff_domain.Pagination) *Pagination {
	_pagination := bsgostuff_types.DerefZero(value)
	_skip := new(int64)
	_limit := new(int64)

	if _pagination.Skip.IsSet() {
		_skip = _pagination.Skip.GetPtr()
	}
	if _pagination.Limit.IsSet() {
		_limit = _pagination.Limit.GetPtr()
	}

	return &Pagination{
		Skip:  _skip,
		Limit: _limit,
	}
}

func FromPaginationProto(value *Pagination) bsgostuff_domain.Pagination {
	_skip := bsgostuff_types.NewInt(0)
	_limit := bsgostuff_types.NewInt(100)

	if value != nil {
		if value.Skip != nil {
			_skip.Set(*value.Skip)
		}
		if value.Limit != nil && *value.Limit != 0 {
			_limit.Set(*value.Limit)
		}
	}

	return bsgostuff_domain.Pagination{
		Skip:  _skip,
		Limit: _limit,
	}
}

func (value *Pagination) FromProto() bsgostuff_domain.Pagination {
	return FromPaginationProto(value)
}
