package main

import (
	"net/http"
	"time"
)

type DbFields struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

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
