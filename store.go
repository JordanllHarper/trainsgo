package main

import "fmt"

const (
	StoreReaderErrIdNotFound storeReaderErrorCode = iota
	StoreReaderErrInternalError
)

const (
	StoreDeleterErrIdNotFound storeDeleterErrorCode = iota
	StoreDeleterErrInternalError
)

type (
	storeReader[T any] interface {
		all() (map[Id]T, *storeReaderError)
		getById(id Id) (T, *storeReaderError)
	}

	storeDeleter interface {
		delete(id Id) *storeDeleterError
	}

	storeReaderDeleter[T any] interface {
		storeReader[T]
		storeDeleter
	}

	storeReaderErrorCode int

	storeReaderError struct {
		id     Id
		entity string
		code   storeReaderErrorCode
	}

	storeDeleterErrorCode int
	storeDeleterError     struct {
		id     Id
		entity string
		code   storeDeleterErrorCode
	}

	storeRenamerErrorCode int
	storeRenamerError     struct {
		id     Id
		entity string
		code   storeRenamerErrorCode
	}
)

func (err storeReaderError) Error() string {
	switch err.code {
	case StoreReaderErrIdNotFound:
		return fmt.Sprintf("%s with ID %v doesn't exist", err.entity, err.id)
	default:
		return fmt.Sprintf("Unrecognised error id %d", err.code)
	}
}

func (err storeDeleterError) Error() string {
	switch err.code {
	case StoreDeleterErrIdNotFound:
		return fmt.Sprintf("%s with ID %v doesn't exist", err.entity, err.id)
	default:
		return fmt.Sprintf("Unrecognised error id %d", err.code)
	}
}

func newStoreReaderError(id Id, entity string, code storeReaderErrorCode) *storeReaderError {
	err := storeReaderError{id, entity, code}
	return &err
}

func newStoreDeleterError(id Id, entity string, code storeDeleterErrorCode) *storeDeleterError {
	err := storeDeleterError{id, entity, code}
	return &err
}
