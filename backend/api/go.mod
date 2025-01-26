module github.com/JordanllHarper/trainsgo/api

go 1.23.4

require github.com/JordanllHarper/trainsgo/app v0.0.0-00010101000000-000000000000

require (
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
)

require (
	golang.org/x/text v0.18.0 // indirect
	gorm.io/driver/mysql v1.5.7
	gorm.io/gorm v1.25.12
)

replace github.com/JordanllHarper/trainsgo/app => ../app
