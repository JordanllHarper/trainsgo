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
		all() (map[id]T, *storeReaderError)
		getById(id id) (T, *storeReaderError)
	}

	storeDeleter interface {
		delete(id id) *storeDeleterError
	}

	storeReaderDeleter[T any] interface {
		storeReader[T]
		storeDeleter
	}

	storeReaderErrorCode int

	storeReaderError struct {
		id     id
		entity string
		code   storeReaderErrorCode
	}

	storeDeleterErrorCode int
	storeDeleterError     struct {
		id     id
		entity string
		code   storeDeleterErrorCode
	}

	storeRenamerErrorCode int
	storeRenamerError     struct {
		id     id
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

func newStoreReaderError(id id, entity string, code storeReaderErrorCode) *storeReaderError {
	err := storeReaderError{id, entity, code}
	return &err
}

func newStoreDeleterError(id id, entity string, code storeDeleterErrorCode) *storeDeleterError {
	err := storeDeleterError{id, entity, code}
	return &err
}
