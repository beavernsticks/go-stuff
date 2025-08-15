package bsgostuff_types

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Enum[T comparable] struct {
	value T
	set   bool
	null  bool
}

func NewEnum[T comparable](v T) Enum[T] {
	return Enum[T]{value: v, set: true}
}

func (e *Enum[T]) Set(value any) error {
	e.set = true
	e.null = false

	switch v := value.(type) {
	case T:
		e.value = v
	case *T:
		if v == nil {
			e.null = true
		} else {
			e.value = *v
		}
	case pgtype.Text:
		if !v.Valid {
			e.null = true
		} else {
			parsed, err := parseEnumFromString[T](v.String)
			if err != nil {
				return err
			}
			e.value = parsed
		}
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

func (e Enum[T]) Get() *T {
	if !e.set || e.null {
		return nil
	}
	return &e.value
}

// parseEnumFromString преобразует строку в значение enum типа T.
// Поддерживает:
//   - Прото-енумы (через protoreflect)
//   - Обычные Go-енумы (int/string)
//   - Регистронезависимые сравнения и snake_case.
func parseEnumFromString[T comparable](s string) (T, error) {
	var zero T
	tType := reflect.TypeOf(zero)

	// 1. Проверяем, является ли T proto-enum (имеет метод Descriptor).
	if protoEnum, ok := any(zero).(interface {
		Descriptor() protoreflect.EnumDescriptor
	}); ok {
		enumDesc := protoEnum.Descriptor()
		// Ищем значение enum по имени (регистронезависимо).
		for i := 0; i < enumDesc.Values().Len(); i++ {
			val := enumDesc.Values().Get(i)
			if strings.EqualFold(string(val.Name()), strings.ToUpper(s)) {
				return reflect.ValueOf(val.Number()).Interface().(T), nil
			}
		}
		return zero, fmt.Errorf("unknown proto enum value: %s", s)
	}

	// 2. Обработка стандартных Go-енумов (int/string).
	switch tType.Kind() {
	case reflect.String:
		// Для строковых enum (type Status string).
		return any(strings.ToUpper(s)).(T), nil

	case reflect.Int, reflect.Int32:
		// Для числовых enum (type Status int).
		for i := 0; i < tType.NumField(); i++ {
			field := tType.Field(i)
			if strings.EqualFold(field.Name, s) {
				return reflect.ValueOf(zero).Field(i).Interface().(T), nil
			}
		}
	}

	return zero, fmt.Errorf("unsupported enum type or unknown value: %s", s)
}
