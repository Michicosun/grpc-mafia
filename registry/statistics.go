package registry

import (
	"net/http"
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

	errors := make([]error, 0)

	for _, report := range round_report.UserReports {
		if err := Server.db.AddNewRound(report.Login, report.Win, report.RoundTime); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) != 0 {
		c.JSON(http.StatusNotFound, gin.H{"errors": errors})
	}

	c.JSON(http.StatusOK, "submitted")
}

func getStatisticsForUser(c *gin.Context) {
	login := c.Param("login")

	stat, err := Server.db.GetStatistics(login)
	if err != nil {
		EndWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, stat)
}

func registerStatisticsRoutes(r *gin.Engine) {
	r.POST("/statistics/round", endRoundHandler)
	r.GET("/statistics/:login", getStatisticsForUser)
}
