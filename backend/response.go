package main

import (
	"net/http"
)

type ResponseBody interface {
	StatusCode() int
}

type HttpError interface {
	Status() (int, string)
}

type ClientError struct {
	message string
}

func (clientError ClientError) Status() (int, string) {
	return http.StatusBadRequest, clientError.message
}

func needBody() ClientError      { return ClientError{"Need body"} }
func provideId() ClientError     { return ClientError{"No ID provided"} }
func invalidId() ClientError     { return ClientError{"Invalid ID"} }
func malformedBody() ClientError { return ClientError{"Malformed body"} }
