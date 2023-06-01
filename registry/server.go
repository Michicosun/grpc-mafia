package registry

import (
	"fmt"
	"grpc-mafia/registry/db"
	"grpc-mafia/registry/fs"

	"github.com/gin-gonic/gin"
)

var Server = &server{}

type server struct {
	router      *gin.Engine
	db          *db.DBAdapter
	ava_storage *fs.FileStorage
	pdf_storage *fs.FileStorage
}

func (s *server) Init(data_folder string) {
	avas, err := fs.CreateFileStorage(fmt.Sprintf("%s/avatars", data_folder))
	if err != nil {
		panic(err)
	}

	pdfs, err := fs.CreateFileStorage(fmt.Sprintf("%s/pdfs", data_folder))
	if err != nil {
		panic(err)
	}

	s.router = gin.Default()

	s.ava_storage = avas
	s.pdf_storage = pdfs

	s.db = db.NewDBAdapter(fmt.Sprintf("%s/db", data_folder))
	Server.db.InitTables()

	registerAuxiliaryRoutes(s.router)
	registerUsersRoutes(s.router)
	// register statistics routes

}

func (s *server) Run() {
	s.router.Run()
}
