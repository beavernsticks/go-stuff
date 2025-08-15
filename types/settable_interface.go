package bsgostuff_types

type Settable interface {
	Set(any) error
	Get() any
	GetValue() any
	GetPtr() any
	IsSet() bool
	IsNull() bool
	ToProto() any
	ToPgx() any
}

func parseString[T any](s string, parser func(string) (T, error)) (T, error) {
	var zero T
	if s == "" {
		return zero, nil
	}

	return parser(s)
}

func toPgxWrapper[T any](value T, valid bool, pgxType interface {
	Valid(bool)
	Set(T)
}) any {
	pgxType.Set(value)
	pgxType.Valid(valid)
	return pgxType
}
