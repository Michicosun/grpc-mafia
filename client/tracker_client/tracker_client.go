package tracker_client

import (
	"context"
	"fmt"
	"grpc-mafia/util"
	"net/http"
	"time"

	"github.com/Khan/genqlient/graphql"
)

var TrackerClient = &trackerClient{}

type trackerClient struct {
	client graphql.Client
}

func (tc *trackerClient) AddComment(round_id string, from string, text string) (*AddCommentResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return AddComment(ctx, tc.client, round_id, from, text)
}

func (tc *trackerClient) GetRoundInfo(round_id string) (*GetRoundInfoResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return GetRoundInfo(ctx, tc.client, round_id)
}

func (tc *trackerClient) ListRounds(n int, state string) (*ListRoundsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return ListRounds(ctx, tc.client, n, RoundState(state))
}

func (tc *trackerClient) Init() {
	host := util.GetEnvWithDefault("TRACKER_HOST", "localhost")
	port := util.GetEnvWithDefault("TRACKER_PORT", "9090")

	endpoint := fmt.Sprintf("http://%s:%s/graphql", host, port)

	tc.client = graphql.NewClient(endpoint, http.DefaultClient)
}
