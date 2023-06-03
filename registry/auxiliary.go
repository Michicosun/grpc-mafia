package registry

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingResponse struct {
	Text string `json:"text"`
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, PingResponse{
		Text: "pong",
	})
}

func registerAuxiliaryRoutes(r *gin.Engine) {
	r.GET("/ping", ping)
}
