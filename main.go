package main

import (
	"log"

	"github.com/polyxia-org/morty-registry/internal/registry"
)

func main() {
	reg, err := registry.NewServer()
	if err != nil {
		log.Fatal(err)
	}
	reg.Serve()
}
