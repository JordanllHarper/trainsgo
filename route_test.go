package main

import (
	"slices"
	"testing"

	"github.com/google/uuid"
)

func Test_routerImpl_Route(t *testing.T) {

	stid1 := uuid.MustParse("5a71e7f1-b037-4e97-a618-bd6fd8043ed0")
	stid2 := uuid.MustParse("5a71e7f1-b037-4e97-a618-bd6fd8043ed1")
	stid3 := uuid.MustParse("5a71e7f1-b037-4e97-a618-bd6fd8043ed2")
	lid1 := uuid.MustParse("5a71e7f1-b037-4e97-a618-bd6fd8043ed3")
	lid2 := uuid.MustParse("5a71e7f1-b037-4e97-a618-bd6fd8043ed4")

	tests := []struct {
		name      string
		currentId Id
		destId    Id
		want      []Line
		want2     error
	}{
		{
			name:      "Route lines with 2 nodes together",
			currentId: stid1,
			destId:    stid2,
			want: []Line{
				{
					Id:         lid1,
					StationOne: stid1,
					StationTwo: stid2,
				},
			},
		},
		{
			name:      "Route lines with 3 nodes connected",
			currentId: stid1,
			destId:    stid3,
			want: []Line{
				{
					Id:         lid1,
					StationOne: stid1,
					StationTwo: stid2,
				},
				{
					Id:         lid2,
					StationOne: stid2,
					StationTwo: stid3,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stations := stationStoreLocal{
				map[Id]Station{
					stid1: {
						E: Entity{
							Id:  stid1,
							Pos: Position{0, 0},
						},
						SurroundingLines: []Line{
							{
								Id:         lid1,
								StationOne: stid1,
								StationTwo: stid2,
							},
						},
					},
					stid2: {
						E: Entity{
							Id:  stid2,
							Pos: Position{10, 10},
						},
						SurroundingLines: []Line{
							{
								Id:         lid1,
								StationOne: stid2,
								StationTwo: stid1,
							},
						},
					},
					stid3: {
						E: Entity{
							Id:  stid3,
							Pos: Position{20, 20},
						},
						SurroundingLines: []Line{
							{
								Id:         lid2,
								StationOne: stid2,
								StationTwo: stid3,
							},
						},
					},
				},
			}
			ri := routeStoreLocal{
				map[Id]Route{},
				stations,
			}
			got, got2 := ri.MapRoute(tt.currentId, tt.destId)
			if !slices.Equal(got, tt.want) {
				t.Errorf("MapRoute() = \n%v,\n want \n%v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("MapRoute() = \n%v,\n want \n%v", got2, tt.want2)
			}
		})
	}
}
