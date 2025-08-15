package bsgostuff_types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

type Bool struct {
	value bool
	set   bool
	null  bool
}

// --- Универсальные методы ---

// Set принимает: bool, *bool, pgtype.Bool, string ("true"/"1"), nil.
func (b *Bool) Set(value any) error {
	b.set = true
	b.null = false
	b.value = false // Сбрасываем значение

	switch v := value.(type) {
	case bool:
		b.value = v
	case *bool:
		if v == nil {
			b.null = true
		} else {
			b.value = *v
		}
	case pgtype.Bool:
		if !v.Valid {
			b.null = true
		} else {
			b.value = v.Bool
		}
	case string:
		return b.parseFromString(v)
	case *string:
		if v == nil {
			b.null = true
		} else {
			return b.parseFromString(*v)
		}
	case nil:
		b.null = true
	default:
		return fmt.Errorf("unsupported bool type: %T", value)
	}
	return nil
}

// Get возвращает bool (или false если null/unset).
func (b Bool) Get() bool {
	if !b.set || b.null {
		return false
	}
	return b.value
}

// GetPtr возвращает *bool (nil если null/unset).
func (b Bool) GetPtr() *bool {
	if !b.set || b.null {
		return nil
	}
	return &b.value
}

// --- Специфичные методы ---

// ToPgxBool конвертирует в pgtype.Bool.
func (b Bool) ToPgxBool() pgtype.Bool {
	if !b.set || b.null {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: b.value, Valid: true}
}

// --- Приватные методы ---

func (b *Bool) parseFromString(s string) error {
	val, err := strconv.ParseBool(s)
	if err != nil {
		return fmt.Errorf("invalid bool format: %w", err)
	}
	b.value = val
	return nil
}

// --- Интеграция с БД и JSON ---

// Scan для sql.Scanner.
func (b *Bool) Scan(src any) error {
	switch v := src.(type) {
	case bool:
		return b.Set(v)
	case []byte:
		return b.parseFromString(string(v))
	case string:
		return b.parseFromString(v)
	case nil:
		b.null = true
		b.set = true
	default:
		return fmt.Errorf("unsupported scan type: %T", src)
	}
	return nil
}

// Value для driver.Valuer.
func (b Bool) Value() (driver.Value, error) {
	if !b.set || b.null {
		return nil, nil
	}
	return b.value, nil
}

// MarshalJSON для JSON.
func (b Bool) MarshalJSON() ([]byte, error) {
	if !b.set || b.null {
		return []byte("null"), nil
	}
	return json.Marshal(b.value)
}

// UnmarshalJSON для JSON.
func (b *Bool) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		b.null = true
		return nil
	}
	var val bool
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	b.value = val
	b.set = true
	return nil
}

func (b Bool) IsNull() bool {
	return b.null
}

func (b Bool) IsSet() bool {
	return b.set
}

func NewBool(value bool) Bool {
	return Bool{value: value, set: true}
}
