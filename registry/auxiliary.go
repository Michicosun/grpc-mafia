package registry

import "github.com/gin-gonic/gin"

type PingResponse struct {
	Text string `json:"text"`
}

func ping(c *gin.Context) {
	c.JSON(200, PingResponse{
		Text: "pong",
	})
}

func registerAuxiliaryRoutes(r *gin.Engine) {
	r.GET("/ping", ping)
}
