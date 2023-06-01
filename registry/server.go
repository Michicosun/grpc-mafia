package registry

import (
	"fmt"
	"grpc-mafia/registry/db"
	"grpc-mafia/registry/fs"
	"grpc-mafia/registry/pdfgen"
	"log"

	"github.com/gin-gonic/gin"
)

var Server = &server{}

type server struct {
	router      *gin.Engine
	db          *db.DBAdapter
	ava_storage *fs.FileStorage
	pdf_storage *fs.FileStorage
	gen         *pdfgen.PdfGen

	render_requests chan pdfgen.RenderRequest
}

func (s *server) WorkerRenderLoop() {
	for request := range s.render_requests {
		pdf_data, err := s.gen.Render(request)

		if err != nil {
			log.Fatal(err)
		}

		s.pdf_storage.Write(request.OutFile, pdf_data)
	}
}

func (s *server) SubmitRenderRequest(request pdfgen.RenderRequest) {
	s.render_requests <- request
}

func (s *server) Init(data_folder string, renders int) {
	s.router = gin.Default()

	avas, err := fs.CreateFileStorage(fmt.Sprintf("%s/avatars", data_folder))
	if err != nil {
		panic(err)
	}
	s.ava_storage = avas

	pdfs, err := fs.CreateFileStorage(fmt.Sprintf("%s/pdfs", data_folder))
	if err != nil {
		panic(err)
	}
	s.pdf_storage = pdfs

	s.db = db.NewDBAdapter(fmt.Sprintf("%s/db", data_folder))
	Server.db.InitTables()

	s.gen = pdfgen.NewPDFGen()
	s.render_requests = make(chan pdfgen.RenderRequest, 10*renders)
	for i := 0; i < renders; i += 1 {
		go s.WorkerRenderLoop()
	}

	registerAuxiliaryRoutes(s.router)
	registerUsersRoutes(s.router)
	// register statistics routes
	registerPdfRoutes(s.router)
}

func (s *server) Run() {
	s.router.Run()
}
