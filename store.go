package main

import "fmt"

const (
	StoreReaderErrIdNotFound StoreReaderErrorCode = iota
	StoreReaderErrInternalError
)

const (
	StoreDeleterErrIdNotFound StoreDeleterErrorCode = iota
	StoreDeleterErrInternalError
)

type (
	StoreReader[T any] interface {
		All() (map[Id]T, *StoreReaderError)
		GetById(id Id) (T, *StoreReaderError)
	}

	StoreDeleter interface {
		Delete(id Id) *StoreDeleterError
	}

	StoreReaderDeleter[T any] interface {
		StoreReader[T]
		StoreDeleter
	}

	StoreReaderErrorCode int

	StoreReaderError struct {
		id     Id
		entity string
		code   StoreReaderErrorCode
	}

	StoreDeleterErrorCode int
	StoreDeleterError     struct {
		id     Id
		entity string
		code   StoreDeleterErrorCode
	}

	StoreRenamerErrorCode int
	StoreRenamerError     struct {
		id     Id
		entity string
		code   StoreRenamerErrorCode
	}
)

func (err StoreReaderError) Error() string {
	switch err.code {
	case StoreReaderErrIdNotFound:
		return fmt.Sprintf("%s with ID %v doesn't exist", err.entity, err.id)
	default:
		return fmt.Sprintf("Unrecognised error id %d", err.code)
	}
}

func (err StoreDeleterError) Error() string {
	switch err.code {
	case StoreDeleterErrIdNotFound:
		return fmt.Sprintf("%s with ID %v doesn't exist", err.entity, err.id)
	default:
		return fmt.Sprintf("Unrecognised error id %d", err.code)
	}
}

func NewStoreReaderError(id Id, entity string, code StoreReaderErrorCode) *StoreReaderError {
	err := StoreReaderError{id, entity, code}
	return &err
}

func NewStoreDeleterError(id Id, entity string, code StoreDeleterErrorCode) *StoreDeleterError {
	err := StoreDeleterError{id, entity, code}
	return &err
}
