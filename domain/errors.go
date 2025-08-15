package bsgostuff_domain

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrDuplicate       = errors.New("duplicate")
	ErrInternal        = errors.New("internal error")
	ErrForbidden       = errors.New("forbidden")
)

// Хелперы для проверки типа ошибки
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func IsForbiddenError(err error) bool {
	return errors.Is(err, ErrForbidden)
}

func IsInvalidArgumentError(err error) bool {
	return errors.Is(err, ErrInvalidArgument)
}

func IsDuplicateError(err error) bool {
	return errors.Is(err, ErrDuplicate)
}

func IsInternaltError(err error) bool {
	return errors.Is(err, ErrInternal)
}
