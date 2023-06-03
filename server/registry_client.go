package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grpc-mafia/registry"
	"grpc-mafia/util"
	"net/http"
	"time"

	zlog "github.com/rs/zerolog/log"
)

var RegistryClient = &registryClient{}

type registryClient struct {
	registryEndpoint string
}

func (rc *registryClient) MakeUserReport(login string, is_win uint8, rt time.Duration) registry.UserRoundReport {
	return registry.UserRoundReport{
		Login:     login,
		Win:       is_win,
		RoundTime: rt,
	}
}

func (rc *registryClient) MakeHandleEndpoint() string {
	return fmt.Sprintf("http://%s/statistics/round", rc.registryEndpoint)
}

func (rc *registryClient) SendRoundReport(g *Game) {
	round_end_timestamp := time.Now()
	round_duration := round_end_timestamp.Sub(g.game_start_timestamp)

	report := registry.RoundReport{
		UserReports: make([]registry.UserRoundReport, 0),
	}

	if g.state == WinMafia {
		for m := range g.mafia {
			report.UserReports = append(report.UserReports, rc.MakeUserReport(m, 1, round_duration))
		}
	} else {
		for m := range g.sheriffs {
			report.UserReports = append(report.UserReports, rc.MakeUserReport(m, 1, round_duration))
		}
		for m := range g.alive_players {
			report.UserReports = append(report.UserReports, rc.MakeUserReport(m, 1, round_duration))
		}
	}

	for m := range g.ghosts {
		report.UserReports = append(report.UserReports, rc.MakeUserReport(m, 0, round_duration))
	}

	content, _ := json.Marshal(report)

	_, err := http.Post(
		rc.MakeHandleEndpoint(),
		"application/json",
		bytes.NewBuffer(content),
	)

	if err != nil {
		zlog.Error().Err(err).Msg("while submitting game report to registry")
	}
}

func (rc *registryClient) Init() {
	host := util.GetEnvWithDefault("REGISTRY_HOST", "localhost")
	port := util.GetEnvWithDefault("REGISTRY_PORT", "8080")

	rc.registryEndpoint = fmt.Sprintf("%s:%s", host, port)

	zlog.Info().Str("endpoint", rc.registryEndpoint).Msg("registry client configured")
}
