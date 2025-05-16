package main

import (
	"testing"

	"github.com/google/uuid"
)

const id1 = "ef8cfc17-21c3-458d-8deb-d96bc423759e"
const id2 = "df652772-cae8-4ce2-9c6e-2bcb10a37762"

var dummyPos position = position{10, 10}

func Test_stationStoreLocal_getById(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		id      id
		want    station
		wantErr bool
	}{
		{
			name: "Gets Station successfully",
			id:   uuid.MustParse(id1),
			want: station{
				entity: entity{
					id: uuid.MustParse(id1),
				},
			},
		},
		{
			name:    "Errors if station doesn't exist",
			id:      uuid.MustParse(id2),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ssl := stationStoreLocal{
				stations: map[id]station{
					uuid.MustParse(id1): {
						entity: entity{
							id: uuid.MustParse(id1),
						},
					},
				},
			}
			got, gotErr := ssl.getById(tt.id)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("getById() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("getById() succeeded unexpectedly")
			}
			if got.id != tt.want.id {
				t.Errorf("getById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stationStoreLocal_register(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		s       station
		wantErr bool
	}{
		{
			name: "Register station successful",
			s: station{
				entity: entity{
					id: uuid.MustParse(id2),
				},
			},
			wantErr: false,
		},
		{
			name: "Register station duplicate id",
			s: station{
				entity: entity{
					id: uuid.MustParse(id1),
				},
			},
			wantErr: true,
		},
		{
			name: "Register station duplicate pos",
			s: station{
				entity: entity{
					id:       uuid.MustParse(id2),
					position: dummyPos,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ssl := stationStoreLocal{
				stations: map[id]station{
					uuid.MustParse(id1): {
						entity: entity{
							id:       uuid.MustParse(id1),
							position: dummyPos,
						},
					},
				},
			}
			gotErr := ssl.register(tt.s)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("register() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("register() succeeded unexpectedly")
			}
		})
	}
}
