package bsgostuff_types

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type String struct {
	value string
	set   bool
	null  bool
}

func NewString(v string) String {
	return String{value: v, set: true}
}

func (s *String) Set(value any) error {
	s.set = true
	s.null = false

	switch v := value.(type) {
	case string:
		s.value = v
	case *string:
		if v == nil {
			s.null = true
		} else {
			s.value = *v
		}
	case pgtype.Text:
		s.null = !v.Valid
		s.value = v.String
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

func (s String) Get() any { return s.GetPtr() }
func (s String) GetValue() string {
	if !s.set || s.null {
		return ""
	}
	return s.value
}
func (s String) GetPtr() *string {
	if !s.set || s.null {
		return nil
	}
	return &s.value
}
func (s String) IsSet() bool        { return s.set }
func (s String) IsNull() bool       { return s.null }
func (s String) ToProto() *string   { return s.GetPtr() }
func (s String) ToPgx() pgtype.Text { return pgtype.Text{String: s.value, Valid: s.set && !s.null} }
