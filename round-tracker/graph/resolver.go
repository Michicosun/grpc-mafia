package graph

import (
	"grpc-mafia/round-tracker/db"
)

//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	db *db.DBAdapter
}

func NewResolver(file string) *Resolver {
	db := db.NewDBAdapter(file)
	db.InitTables()

	return &Resolver{
		db: db,
	}
}
