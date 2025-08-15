package bsgostuff_domain

import bsgostuff_types "github.com/beavernsticks/go-stuff/types"

type Sort struct {
	Field      bsgostuff_types.String
	IsReversed bsgostuff_types.Bool
}
