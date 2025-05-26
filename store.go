package main

// Atomic interfaces
type (
	Store[T any] interface {
		All() (map[Id]T, error)
	}

	IDable[T any] interface {
		GetById(id Id) (T, error)
	}

	Nameable[T any] interface {
		GetByName(name string) ([]T, error)
	}

	Deleter interface {
		Delete(id Id) error
		DeleteBatch(id []Id) error
	}
)

// Composed interfaces
type (
	StoreIDable[T any] interface {
		Store[T]
		IDable[T]
	}

	StoreIdNameable[T any] interface {
		StoreIDable[T]
		Nameable[T]
	}

	StoreDeleter[T any] interface {
		Store[T]
		Deleter
	}
)
