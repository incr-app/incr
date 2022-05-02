package main

import (
	"fmt"
	"log"
	"strings"

	"incr.app/incr/common"
	"incr.app/incr/datastore"
	"incr.app/incr/datastore/json"
	"incr.app/incr/schema"
)

type App struct {
	datastore.Datastore
	schema.Registry
}

func main() {
	// This is all placeholder code for now. Eventually we'll be querying the server to ask for all this
	myDatastore, err := json.NewDatastore("versions.json")
	if err != nil {
		panic(err)
	}
	app := &App{
		Datastore: myDatastore,
		Registry:  schema.DefaultRegistry,
	}
	registry := app.Registry
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

	// id :=

	conn, err := app.Datastore.Connect()
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	var id datastore.Id = "incr"
	initial, err := conn.Get(id)
	if err != nil {
		if datastore.IsVersionNotFoundError(err) {
			log.Printf("Could not find id: %s. Using 1.2.3\n", id)
			initial = semverSchema.New("1.2.3")
		}
	}

	fmt.Printf("New semver version created: %s\n", initial.String())

	bumpParams := common.BumpParams{
		"type": "major",
	}
	afterBump, err := semverSchema.Bump(initial, bumpParams)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Bumped %s with {\"type\": \"major\"}, got %s\n", initial, afterBump)

	err = conn.Put(id, afterBump)
	if err != nil {
		log.Fatal(err)
	}
}
