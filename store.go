package main

import "net/http"

const (
	StoreErrorIdNotFound    StoreErrorCode = 0
	StoreErrorInternalError StoreErrorCode = 1
)

type StoreErrorCode int
type StoreError interface {
	StoreErrorCode() StoreErrorCode
	HttpError
}

type (
	StoreReader[T any] interface {
		All() (map[Id]T, StoreError)
		GetById(id Id) (T, StoreError)
	}

	StoreDeleter interface {
		Delete(id Id) StoreError
	}

	StoreReaderDeleter[T any] interface {
		StoreReader[T]
		StoreDeleter
	}

	IdDoesntExist Id
	InternalError struct{ error }
)

func (e IdDoesntExist) StoreErrorCode() StoreErrorCode { return StoreErrorIdNotFound }
func (e IdDoesntExist) HttpCode() int                  { return http.StatusBadRequest }
func (e IdDoesntExist) Error() string                  { return msgIdDoesntExist(Id(e)) }

func (e InternalError) HttpCode() int                  { return http.StatusInternalServerError }
func (e InternalError) StoreErrorCode() StoreErrorCode { return StoreErrorIdNotFound }
func (e InternalError) Error() string                  { return e.error.Error() }
