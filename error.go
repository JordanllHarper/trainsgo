package main

import "fmt"

func errorIdNotFound(id id, entity string) error {
	return fmt.Errorf("%s with ID %v doesn't exist", entity, id)
}

func idAlreadyExists(id id, entity string) error {
	return fmt.Errorf(
		"%v with ID %s already exists",
		entity,
		id,
	)
}
