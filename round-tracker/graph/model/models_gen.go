// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Comment struct {
	From string `json:"from"`
	Text string `json:"text"`
}

type Player struct {
	Login string `json:"login"`
	Role  string `json:"role"`
	Alive bool   `json:"alive"`
}

type PlayerInfo struct {
	Login string `json:"login"`
	Role  string `json:"role"`
}

type PlayerStatus struct {
	Login string `json:"login"`
	Alive bool   `json:"alive"`
}

type Round struct {
	ID        string     `json:"id"`
	State     RoundState `json:"state"`
	StartedAt string     `json:"started_at"`
	Players   []*Player  `json:"players"`
	Comments  []*Comment `json:"comments"`
}

type RoundInfo struct {
	ID        string     `json:"id"`
	State     RoundState `json:"state"`
	StartedAt string     `json:"started_at"`
}

type RoundState string

const (
	RoundStateRunning  RoundState = "RUNNING"
	RoundStateFinished RoundState = "FINISHED"
)

var AllRoundState = []RoundState{
	RoundStateRunning,
	RoundStateFinished,
}

func (e RoundState) IsValid() bool {
	switch e {
	case RoundStateRunning, RoundStateFinished:
		return true
	}
	return false
}

func (e RoundState) String() string {
	return string(e)
}

func (e *RoundState) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RoundState(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RoundState", str)
	}
	return nil
}

func (e RoundState) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
