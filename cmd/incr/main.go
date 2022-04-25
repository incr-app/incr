package main

import (
	"fmt"
	"strings"

	"incr.app/incr/common"
	"incr.app/incr/schema"
)

func main() {
	// This is all placeholder code for now. Eventually we'll be querying the server to ask for all this
	registry := schema.DefaultRegistry
	embeddedSchemas, err := schema.GetEmbeddedSchemas()
	if err != nil {
		panic(err)
	}
	for schemaName, schema := range embeddedSchemas {
		registry.Register(schemaName, schema)
	}

	fmt.Printf("Registry has %d registered schemas: %s\n", len(registry.Schemas()), strings.Join(registry.Schemas(), ", "))

	semverSchema, _ := registry.Lookup("semver")
	fmt.Printf("Loaded %s\n", semverSchema.Name())
	initial := semverSchema.New("1.2.3")

	fmt.Printf("New semver version created: %s\n", initial.String())

	bumpParams := common.BumpParams{
		"type": "major",
	}
	afterBump, err := semverSchema.Bump(initial, bumpParams)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Bumped %s with {\"type\": \"major\"}, got %s\n", initial, afterBump)
}
