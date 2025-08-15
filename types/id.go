package bsgostuff_types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ID struct {
	value uuid.UUID
	set   bool
	null  bool
}

// --- Универсальные методы ---

// Set принимает: string, uuid.UUID, pgtype.UUID, *string, nil.
func (u *ID) Set(value any) error {
	u.set = true
	u.null = false
	u.value = uuid.Nil // Сбрасываем значение

	switch v := value.(type) {
	case string:
		return u.parseFromString(v)
	case *string:
		if v == nil {
			u.null = true
		} else {
			return u.parseFromString(*v)
		}
	case uuid.UUID:
		u.value = v
	case pgtype.UUID:
		if !v.Valid {
			u.null = true
		} else {
			u.value = v.Bytes
		}
	case nil:
		u.null = true
	default:
		return fmt.Errorf("unsupported UUID type: %T", value)
	}
	return nil
}

// Get возвращает uuid.UUID (или uuid.Nil если null/unset).
func (u ID) Get() uuid.UUID {
	if !u.set || u.null {
		return uuid.Nil
	}
	return u.value
}

// GetPtr возвращает *uuid.UUID (nil если null/unset).
func (u ID) GetPtr() *uuid.UUID {
	if !u.set || u.null {
		return nil
	}
	return &u.value
}

// --- Специфичные методы ---

// NewUUID генерирует новый UUID и устанавливает его.
func (u *ID) NewID() {
	u.value = uuid.New()
	u.set = true
	u.null = false
}

// ToPgx конвертирует в pgtype.UUID.
func (u ID) ToPgx() pgtype.UUID {
	if !u.set || u.null {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: u.value, Valid: true}
}

// --- Приватные методы ---

func (u *ID) parseFromString(s string) error {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return fmt.Errorf("invalid UUID format: %w", err)
	}
	u.value = parsed
	return nil
}

// --- Интеграция с БД и JSON ---

// Scan для sql.Scanner.
func (u *ID) Scan(src any) error {
	switch v := src.(type) {
	case []byte:
		return u.Set(string(v))
	case string:
		return u.Set(v)
	case nil:
		u.null = true
		u.set = true
	default:
		return fmt.Errorf("unsupported scan type: %T", src)
	}
	return nil
}

// Value для driver.Valuer.
func (u ID) Value() (driver.Value, error) {
	if !u.set || u.null {
		return nil, nil
	}
	return u.value.String(), nil
}

// MarshalJSON для JSON.
func (u ID) MarshalJSON() ([]byte, error) {
	if !u.set || u.null {
		return []byte("null"), nil
	}
	return []byte(`"` + u.value.String() + `"`), nil
}

// UnmarshalJSON для JSON.
func (u *ID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		u.null = true
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	return u.Set(s)
}

func (u ID) IsNull() bool {
	return u.null
}

func (u ID) IsSet() bool {
	return u.set
}

func (u ID) Validate() error {
	if u.set && !u.null && u.value == uuid.Nil {
		return errors.New("UUID не может быть нулевым")
	}
	return nil
}

func NewID() ID {
	u := ID{}
	u.NewID()
	return u
}
