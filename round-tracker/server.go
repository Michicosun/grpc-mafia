package main

import (
	"grpc-mafia/round-tracker/graph"
	"grpc-mafia/util"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const DEFAULT_TRACKER_PORT = "8080"

func main() {
	port := util.GetEnvWithDefault("PORT", DEFAULT_TRACKER_PORT)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver()}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
