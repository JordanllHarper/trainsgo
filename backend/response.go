package main

import (
	"time"
)

type DbFields struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type ResponseBody interface{}

type HttpError interface {
	Status() (int, string)
}

type ClientError struct {
	message string
}

type ServerError struct {
	message string
}


func (serverError ServerError) Status() (int, string) {
	return 500, serverError.message
}

func (clientError ClientError) Status() (int, string) {
	return 400, clientError.message
}
