package main

import (
	"testing"

	"github.com/google/uuid"
)

func Test_routerImpl_Route(t *testing.T) {

	routeId := uuid.MustParse("98e895a2-2fed-40b4-a681-ff6dac7d54a4")
	stid1 := uuid.MustParse("06e83552-27ba-4ab3-a817-4d6538275640")
	stid2 := uuid.MustParse("5a71e7f1-b037-4e97-a618-bd6fd8043ed1")
	lid1 := uuid.MustParse("5a71e7f1-b037-4e97-a618-bd6fd8043ed2") // difference by last

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		currentId Id
		destId    Id
		want      Route
		want2     error
	}{
		{
			name:      "Route with 2 nodes together",
			currentId: stid1,
			destId:    stid2,
			want: Route{
				Id:        routeId,
				CurrentId: stid1,
				DestId:    stid2,
				NextId:    lid1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stations := stationStoreLocal{
				map[Id]Station{
					stid1: Station{},
				},
			}
			ri := routeStoreLocal{
				map[Id]Route{},
				stations,
			}
			got, got2 := ri.MapRoute(tt.currentId, tt.destId)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Route() = %v, want %v", got, tt.want)
			}
			if true {
				t.Errorf("Route() = %v, want %v", got2, tt.want2)
			}
		})
	}
}
