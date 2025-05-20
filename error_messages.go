package main

import "fmt"

func msgMalformedBody() string { return "Malformed body" }

func msgBadId(id string) string           { return fmt.Sprintf("Bad ID %s", id) }
func msgIdDoesntExist(id Id) string       { return fmt.Sprintf("ID %s doesnt exist", id) }
func msgIdAlreadyExists(id Id) string     { return fmt.Sprintf("ID %s already exists", id) }
func msgMethodNotAllowed(m string) string { return fmt.Sprintf("METHOD %s not allowed", m) }
