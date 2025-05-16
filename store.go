package main

type (

	/*
		Represents something

		Use all() to get all entities in the storeReader or getById() for a specific entity.

		Use getByName() to get a list of T with the appropriate name.
	*/
	storeReader[T any] interface {
		all() (map[id]T, error)
		getById(id id) (T, error)
		getByName(name string) ([]T, error)
	}

	/*
		Represents something that holds a storeWriter of T.

		To register() means the network needs to know of a new entity, and will accomodate that.

		deregister() will not immediately delete an entity.
		It will wait for the entity to finish it's tasks before removing from the network.
	*/
	storeWriter[T any] interface {
		// Register a T with the store.
		register(s T) error

		// Deregister a T with the store by ID.
		deregister(id id) error

		// Change the name of a given store item.
		changeName(id id, newName string) error
	}

	storeReaderWriter[T any] interface {
		storeReader[T]
		storeWriter[T]
	}
)
