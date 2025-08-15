package bsgostuff_domain

import (
	bsgostuff_types "github.com/beavernsticks/go-stuff/types"
)

type Pagination struct {
	Skip  bsgostuff_types.Int
	Limit bsgostuff_types.Int
}
