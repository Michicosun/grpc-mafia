package main

import (
	"fmt"
	"grpc-mafia/round-tracker/graph"
	"grpc-mafia/util"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const (
	DEFAULT_TRACKER_PORT = "9090"
	DEFAULT_DATA_FOLDER  = "/tmp/grpc-mafia"
)

func main() {
	port := util.GetEnvWithDefault("PORT", DEFAULT_TRACKER_PORT)
	folder := util.GetEnvWithDefault("DATA_FOLDER", DEFAULT_DATA_FOLDER)

	if err := util.CreateIfNotExists(folder); err != nil {
		panic(err)
	}

	db_file := fmt.Sprintf("%s/tracker-db", folder)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(db_file)}))

	http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	http.Handle("/graphql", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
