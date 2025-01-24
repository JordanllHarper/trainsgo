package main

import (
	"fmt"
	"strings"
	"time"
)

type DbFields struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func stringEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}
