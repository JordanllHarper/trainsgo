package main

type (
	StoreReader[T any] interface {
		All() (map[Id]T, error)
		GetById(id Id) (T, error)
	}

	StoreDeleter interface {
		Delete(id Id) error
		DeleteBatch(id []Id) error
	}

	StoreReaderDeleter[T any] interface {
		StoreReader[T]
		StoreDeleter
	}

	InternalStoreError struct{ error }
)
