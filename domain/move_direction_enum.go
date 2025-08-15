package bsgostuff_domain

type MoveDirectionEnum string

const (
	MoveDirectionEnumUnknown MoveDirectionEnum = ""
	MoveDirectionEnumUp      MoveDirectionEnum = "UP"
	MoveDirectionEnumDown    MoveDirectionEnum = "DOWN"
	MoveDirectionEnumTop     MoveDirectionEnum = "TOP"
	MoveDirectionEnumBottom  MoveDirectionEnum = "BOTTOM"
)

func (s MoveDirectionEnum) Valid() bool {
	switch s {
	case MoveDirectionEnumUp, MoveDirectionEnumDown, MoveDirectionEnumTop, MoveDirectionEnumBottom:
		return true
	default:
		return false
	}
}
