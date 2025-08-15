package bsgostuff_types

type Strings = Slice[string]

func NewStrings(v []string) Strings {
	return Strings{values: v, set: true}
}
