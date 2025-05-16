package main

import "fmt"

type errIdNotFound struct {
	id     id
	entity string
}

type errIdAlreadyExists struct {
	id     id
	entity string
}

func (err errIdNotFound) Error() string {
	return fmt.Sprintf("%s with ID %v doesn't exist", err.entity, err.id)
}

func (err errIdAlreadyExists) Error() string {
	return fmt.Sprintf(
		"%v with ID %s already exists",
		err.entity,
		err.id,
	)
}

func newErrIdNotFound(id id, entity string) error {
	return errIdNotFound{id, entity}
}

func newErrIdAlreadyExists(id id, entity string) error {
	return errIdAlreadyExists{id, entity}
}
