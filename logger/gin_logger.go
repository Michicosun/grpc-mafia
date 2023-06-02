package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		path := c.Request.URL.Path

		c.Next()

		data := &requestData{
			status:    c.Writer.Status(),
			client_ip: c.ClientIP(),
			path:      path,
			method:    c.Request.Method,
			latency:   time.Since(t),
			err:       c.Errors.Last(),
		}

		data.log()
	}
}

type requestData struct {
	status    int
	client_ip string
	path      string
	method    string
	latency   time.Duration
	err       *gin.Error
}

func (r *requestData) log() {
	var loggerEvent *zerolog.Event
	var msg string

	switch {
	case 500 <= r.status:
		loggerEvent = zlog.Error().Err(r.err)
		msg = "Internal error"
	case 400 <= r.status:
		loggerEvent = zlog.Warn().Err(r.err)
		msg = "Client Error"
	default:
		loggerEvent = zlog.Info()
		msg = "Success"
	}

	loggerEvent.
		Str("ip", r.client_ip).Str("method", r.method).
		Str("path", r.path).Int("status", r.status).
		Str("latency", r.latency.String()).
		Msg(msg)
}
