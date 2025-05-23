package main

import (
	"errors"
	"net/http"
)

type (
	badId string

	idDoesntExist   Id
	idAlreadyExists Id

	internalError struct{ error }
)

type (
	HttpError interface {
		HttpCode() int
		error
	}

	badRequest       struct{ error }
	methodNotAllowed string
	malformedBody    struct{}
)

func mapToHttpErr(err error) HttpError {
	var rec HttpError
	if errors.As(err, &rec) {
		return rec
	}
	return internalError{err}
}

func (e malformedBody) Error() string { return msgMalformedBody() }
func (e malformedBody) HttpCode() int { return http.StatusBadRequest }

func (e badId) Error() string { return msgBadId(string(e)) }
func (e badId) HttpCode() int { return http.StatusBadRequest }

func (e idDoesntExist) Error() string { return msgIdDoesntExist(Id(e)) }
func (e idDoesntExist) HttpCode() int { return http.StatusBadRequest }

func (e idAlreadyExists) Error() string { return msgIdAlreadyExists(Id(e)) }
func (e idAlreadyExists) HttpCode() int { return http.StatusBadRequest }

func (e methodNotAllowed) Error() string { return msgMethodNotAllowed(string(e)) }
func (e methodNotAllowed) HttpCode() int { return http.StatusMethodNotAllowed }

func (e internalError) Error() string { return e.error.Error() }
func (e internalError) HttpCode() int { return http.StatusInternalServerError }

func (e badRequest) Error() string { return e.error.Error() }
func (e badRequest) HttpCode() int { return http.StatusBadRequest }
