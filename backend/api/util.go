package api

import (
	"fmt"
	"strings"
)

func checkError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func stringEmpty(s string) bool { return strings.TrimSpace(s) == "" }
