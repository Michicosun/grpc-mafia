package registry

import (
	"encoding/json"
	"fmt"
	"grpc-mafia/logger"
	"grpc-mafia/registry/db"
	"grpc-mafia/registry/fs"
	"grpc-mafia/registry/pdfgen"
	"grpc-mafia/registry/queue"

	"github.com/gin-gonic/gin"
	zlog "github.com/rs/zerolog/log"
)

var Server = &server{}

type ServerConfig struct {
	DataFolder  string
	Renders     int
	RabbitCreds queue.RabbitCredentials
}

type server struct {
	router      *gin.Engine
	db          *db.DBAdapter
	ava_storage *fs.FileStorage
	pdf_storage *fs.FileStorage

	gen        *pdfgen.PdfGen
	task_queue *queue.TaskQueue
}

func (s *server) WorkerRenderLoop() {
	for raw_request := range s.task_queue.GetTaskChan() {
		request := pdfgen.RenderRequest{}

		if err := json.Unmarshal(raw_request.Body, &request); err != nil {
			zlog.Error().Err(err).Msg("unmarshal request")
			continue
		}

		pdf_data, err := s.gen.Render(request)

		if err != nil {
			zlog.Error().Err(err).Msg("while render")
			continue
		}

		s.pdf_storage.Write(request.OutFile, pdf_data)
	}
}

func (s *server) SubmitRenderRequest(request pdfgen.RenderRequest) {
	s.task_queue.SubmitTask(request)
}

func (s *server) Init(cfg ServerConfig) {
	s.router = gin.New()
	s.router.Use(logger.GinLogger())

	avas, err := fs.CreateFileStorage(fmt.Sprintf("%s/avatars", cfg.DataFolder))
	if err != nil {
		panic(err)
	}
	s.ava_storage = avas

	pdfs, err := fs.CreateFileStorage(fmt.Sprintf("%s/pdfs", cfg.DataFolder))
	if err != nil {
		panic(err)
	}
	s.pdf_storage = pdfs

	s.db = db.NewDBAdapter(fmt.Sprintf("%s/db", cfg.DataFolder))
	Server.db.InitTables()

	s.gen = pdfgen.NewPDFGen()
	s.task_queue = queue.NewTaskQueue(cfg.RabbitCreds)
	for i := 0; i < cfg.Renders; i += 1 {
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
