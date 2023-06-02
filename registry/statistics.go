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

type RoundReport struct {
	UserReports []UserRoundReport `json:"user_reports"`
}

func endRoundHandler(c *gin.Context) {
	round_report := RoundReport{}

	if err := c.BindJSON(&round_report); err != nil {
		EndWithError(c, err)
		return
	}

	zlog.Info().Interface("round_report", round_report).Msg("parsed reports")

	for _, report := range round_report.UserReports {
		if err := Server.db.AddNewRound(report.Login, report.Win, report.RoundTime); err != nil {
			EndWithError(c, err)
			return
		}
	}

	c.JSON(200, "submitted")
}

func getStatisticsForUser(c *gin.Context) {
	login := c.Param("login")

	stat, err := Server.db.GetStatistics(login)
	if err != nil {
		EndWithError(c, err)
		return
	}

	c.JSON(200, stat)
}

func registerStatisticsRoutes(r *gin.Engine) {
	r.POST("/statistics/round", endRoundHandler)
	r.GET("/statistics/:login", getStatisticsForUser)
}
