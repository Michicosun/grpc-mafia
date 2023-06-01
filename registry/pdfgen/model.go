package pdfgen

import "grpc-mafia/registry/db"

type RenderRequest struct {
	OutFile string
	User    db.User
	// Statistics
}
