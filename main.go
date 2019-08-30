package main

import (
	"context"
	"fmt"
	"time"

	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/introspection"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
	"github.com/samsarahq/thunder/reactive"
)

type post struct {
	Title     string
	Body      string
	CreatedAt time.Time
}

//* Graphql server struct
type server struct {
	posts []post
}

//* Registers root query type
func (s *server) registerQuery(schema *schemabuilder.Schema) {
	obj := schema.Query()

	obj.FieldFunc("posts", func() []post {
		return s.posts
	})
}

//* Registers root mutation type
func (s *server) registerMutation(schema *schemabuilder.Schema) {
	obj := schema.Mutation()
	obj.FieldFunc("echo", func(args struct{ Message string }) string {
		return args.Message
	})
}

//* Registers post type
func (s *server) registerPost(schema *schemabuilder.Schema) {
	obj := schema.Object("Post", post{})
	obj.FieldFunc("age", func(ctx context.Context, p *post) string {
		reactive.InvalidateAfter(ctx, 5*time.Second)
		return time.Since(p.CreatedAt).String()
	})
}

//* Builds graphql schema
func (s *server) schema() *graphql.Schema {
	builder := schemabuilder.NewSchema()
	s.registerQuery(builder)
	s.registerMutation(builder)
	s.registerPost(builder)
	return builder.MustBuild()
}

func main() {
	/*
		//* Instantiates and builds a server, serves schema on port 3030
		server := &server{
			posts: []post{
				{Title: "first post", Body: "testing", CreatedAt: time.Now()},
				{Title: "graphql", Body: "did you hear about Thunder?", CreatedAt: time.Now()},
			},
		}

		schema := server.schema()
		introspection.AddIntrospectionToSchema(schema)

		//* Expose GraphQL POST endpoint.
		http.Handle("/graphql", graphql.HTTPHandler(schema))
		http.Handle("/graphiql/", http.StripPrefix("/graphiql/", graphiql.Handler()))
		fmt.Println("GraphQL Gateway listening on port :3030")
		http.ListenAndServe(":3030", nil)
	*/
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
