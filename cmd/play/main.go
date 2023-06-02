package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grpc-mafia/registry"
	"net/http"
	"time"
)

func main() {
	report := registry.RoundReport{
		UserReports: []registry.UserRoundReport{
			{
				Login:     "a",
				Win:       0,
				RoundTime: 10 * time.Second,
			},
			{
				Login:     "michicosun",
				Win:       1,
				RoundTime: 10000 * time.Second,
			},
			{
				Login:     "b",
				Win:       1,
				RoundTime: 90 * time.Second,
			},
		},
	}

	content, _ := json.Marshal(report)

	resp, err := http.Post(
		"http://localhost:8080/statistics/round",
		"application/json",
		bytes.NewBuffer(content),
	)

	fmt.Println(resp, err)
}
