package main

import (
	"net/http"
)

type (
	HttpError interface {
		HttpCode() int
		error
	}

	genericError struct {
		error
		code int
	}

	badRequest struct{ error }

	methodNotAllowed    string
	malformedBody       struct{}
	internalServerError struct{ error }

	badId string

	idDoesntExist   Id
	idAlreadyExists Id
)

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

func (e internalServerError) Error() string { return e.error.Error() }
func (e internalServerError) HttpCode() int { return http.StatusInternalServerError }

func (e badRequest) Error() string { return e.error.Error() }
func (e badRequest) HttpCode() int { return http.StatusBadRequest }

func (e genericError) HttpCode() int { return e.code }
