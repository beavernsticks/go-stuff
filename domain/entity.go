package bsgostuff_domain

import (
	bsgostuff_types "github.com/beavernsticks/go-stuff/types"
)

type Entity struct {
	ID        bsgostuff_types.ID
	CreatedAt bsgostuff_types.Timestamp
	UpdatedAt bsgostuff_types.Timestamp
}

type DeletableEntity struct {
	ID        bsgostuff_types.ID
	CreatedAt bsgostuff_types.Timestamp
	UpdatedAt bsgostuff_types.Timestamp
	IsDeleted bsgostuff_types.Bool
}

type UnmodifiedEntity struct {
	ID        bsgostuff_types.ID
	CreatedAt bsgostuff_types.Timestamp
}

func NewEntity() Entity {
	return Entity{
		ID:        bsgostuff_types.NewID(),
		CreatedAt: bsgostuff_types.Now(),
		UpdatedAt: bsgostuff_types.Now(),
	}
}

func NewEntityPtr() *Entity {
	entity := NewEntity()
	return &entity
}

func NewDeletableEntity() DeletableEntity {
	return DeletableEntity{
		ID:        bsgostuff_types.NewID(),
		CreatedAt: bsgostuff_types.Now(),
		UpdatedAt: bsgostuff_types.Now(),
		IsDeleted: bsgostuff_types.NewBool(false),
	}
}

func NewDeletableEntityPtr() *DeletableEntity {
	entity := NewDeletableEntity()
	return &entity
}

func NewUnmodifiedEntity() UnmodifiedEntity {
	return UnmodifiedEntity{
		ID:        bsgostuff_types.NewID(),
		CreatedAt: bsgostuff_types.Now(),
	}
}

func NewUnmodifiedEntityPtr() *UnmodifiedEntity {
	entity := NewUnmodifiedEntity()
	return &entity
}
