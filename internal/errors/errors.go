package errors

import (
	"errors"
	"fmt"
	"strings"
)

// Error list.
var (
	ErrInternalDB           = errors.New("internal database error")
	ErrInternalCache        = errors.New("internal cache error")
	ErrInternalServer       = errors.New("internal server error")
	ErrInvalidDBFormat      = errors.New("invalid db address")
	ErrInvalidRequestFormat = errors.New("invalid request format")
	ErrInvalidMessageType   = errors.New("invalid message type")
	ErrInvalidAnimeID       = errors.New("invalid anime id")
	ErrAnimeNotFound        = errors.New("anime not found")
	ErrDataStillNew         = errors.New("data is still new")
)

// ErrRequiredField is error for missing field.
func ErrRequiredField(str string) error {
	return fmt.Errorf("required field %s", str)
}

// ErrGTField is error for greater than field.
func ErrGTField(str, value string) error {
	return fmt.Errorf("field %s must be greater than %s", str, value)
}

// ErrGTEField is error for greater than or equal field.
func ErrGTEField(str, value string) error {
	return fmt.Errorf("field %s must be greater than or equal %s", str, value)
}

// ErrLTField is error for lower than field.
func ErrLTField(str, value string) error {
	return fmt.Errorf("field %s must be lower than %s", str, value)
}

// ErrLTEField is error for lower than or equal field.
func ErrLTEField(str, value string) error {
	return fmt.Errorf("field %s must be lower than or equal %s", str, value)
}

// ErrOneOfField is error for oneof field.
func ErrOneOfField(str, value string) error {
	return fmt.Errorf("field %s must be one of %s", str, strings.Join(strings.Split(value, " "), "/"))
}

// ErrDatetimeField is error for datetime field.
func ErrDatetimeField(str, value string) error {
	return fmt.Errorf("field %s must be in %s format", str, value)
}
