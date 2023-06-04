package graph

import "grpc-mafia/round-tracker/graph/model"

//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	rounds map[string]*model.Round
}

func NewResolver() *Resolver {
	return &Resolver{
		rounds: make(map[string]*model.Round),
	}
}
