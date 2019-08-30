package main

import (
	"fmt"
	"time"

	"github.com/samsarahq/thunder/graphql/introspection"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
)

func run() {
	server := &server{
		posts: []post{
			{Title: "first post", Body: "testing", CreatedAt: time.Now()},
			{Title: "graphql", Body: "did you hear about Thunder?", CreatedAt: time.Now()},
		},
	}

	builderSchema := schemabuilder.NewSchema()
	server.registerQuery(builderSchema)
	server.registerMutation(builderSchema)

	valueJSON, err := introspection.ComputeSchemaJSON(*builderSchema)
	if err != nil {
		panic(err)
	}

	fmt.Print(string(valueJSON))
}
