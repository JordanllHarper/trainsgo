package main

type (
	Queue[T any] interface {
		Push(T)
		Pop() T
	}

	queue[T any] []T
)

func NewQueue[T any]() Queue[T] {
	return queue[T]{}
}

func (q queue[T]) Push(item T) {
	q = append(q, item)
}

func (impl queue[T]) Pop() T {
	item := impl[0]
	impl = impl[1:]
	return item
}
