package main

import (
	"github.com/polyxia-org/morty-registry/internal/registry"
	log "github.com/sirupsen/logrus"
)

func main() {
	// log.SetLevel(log.TraceLevel)
	reg, err := registry.NewServer()
	if err != nil {
		log.Fatalf("failed to initialize the registry: %v", err)
	}

	reg.Serve()
}
