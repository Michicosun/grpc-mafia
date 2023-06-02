package registry

import (
	"fmt"
	"grpc-mafia/registry/pdfgen"
	"net/http"

	"github.com/gin-gonic/gin"
	zlog "github.com/rs/zerolog/log"
)

func createPdf(c *gin.Context) {
	login := c.Param("login")

	user, err := Server.db.GetUser(login)
	if err != nil {
		zlog.Error().Err(err).Msg("get user")
		EndWithError(c, err)
		return
	}

	statistics, err := Server.db.GetStatistics(login)
	if err != nil {
		zlog.Error().Err(err).Msg("get statistics")
		EndWithError(c, err)
		return
	}

	out_file := GenRandomName()

	if len(user.AvatarFilename) != 0 {
		user.AvatarFilename = Server.ava_storage.Pwd(user.AvatarFilename)
	}

	Server.SubmitRenderRequest(pdfgen.RenderRequest{
		OutFile:    out_file,
		User:       *user,
		Statistics: *statistics,
	})

	c.JSON(http.StatusOK, gin.H{"url": fmt.Sprintf("http://localhost:8080/pdf/%s", out_file)}) // TODO: remove localhost
}

func getPdf(c *gin.Context) {
	filename := c.Param("filename")

	if err := Server.pdf_storage.RunStat(filename); err != nil {
		zlog.Error().Err(err).Msgf("file %s not found", filename)
		EndWithError(c, err)
		return
	}

	c.File(Server.pdf_storage.Pwd(filename))
}

func registerPdfRoutes(r *gin.Engine) {
	r.POST("/pdf/:login", createPdf)
	r.GET("/pdf/:filename", getPdf)
}
