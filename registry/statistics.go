package registry

import (
	"time"

	"github.com/gin-gonic/gin"
	zlog "github.com/rs/zerolog/log"
)

type UserStatistics struct {
	SessionsCnt uint64
	WinCnt      uint64
	LoseCnt     uint64
	TotalTime   time.Time
}

type UserRoundReport struct {
	Login     string        `json:"login"`
	Win       uint8         `json:"win"`
	RoundTime time.Duration `json:"round_time"`
}

func endRoundHandler(c *gin.Context) {
	report := UserRoundReport{}

	if err := c.BindJSON(&report); err != nil {
		zlog.Error().Err(err).Msg("bind")
		EndWithError(c, err)
		return
	}

	if err := Server.db.AddNewRound(report.Login, report.Win, report.RoundTime); err != nil {
		zlog.Error().Err(err).Msg("insert into database failed")
		EndWithError(c, err)
		return
	}

	c.JSON(200, "submitted")
}

func getStatisticsForUser(c *gin.Context) {
	login := c.Param("login")

	stat, err := Server.db.GetStatistics(login)
	if err != nil {
		zlog.Error().Err(err).Msg("get statistics")
		EndWithError(c, err)
		return
	}

	c.JSON(200, stat)
}

func registerStatisticsRoutes(r *gin.Engine) {
	r.POST("/statistics/round", endRoundHandler)
	r.GET("/statistics/:login", getStatisticsForUser)
}
