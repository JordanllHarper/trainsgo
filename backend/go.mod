module github.com/JordanllHarper/trainsgo/backend/testing_frontend

go 1.23.5

replace github.com/JordanllHarper/trainsgo/backend/engine => ./engine

replace github.com/JordanllHarper/trainsgo/backend/common => ./common

require (
	github.com/JordanllHarper/trainsgo/backend/common v0.0.0
	github.com/JordanllHarper/trainsgo/backend/engine v0.0.0-20250202104413-ee59937f64bf
)

replace github.com/JordanllHarper/backend/common => ./common
