package main

type (
	/*
		Represents a path from one ID to another using the appropriate lineId.

		Generic due to the ability to add more than just stations in the future.

		Given a train is at currentId, the train needs to then be routed to nextId using lineId.

		Note no inclusion of any kind of trip. Given computing a route might be expensive, we should cache the graph, and invalidate and recompute only when we add a new station or line.
	*/
	Route struct {
		CurrentId Id
		NextId    Id
		LineId    Id
	}
)

func NewRoute(
	currentId Id,
	nextId Id,
	lineId Id,
) Route {
	return Route{
		currentId,
		nextId,
		lineId,
	}
}
