package db

import (
	"fmt"
	"grpc-mafia/round-tracker/graph/model"
	"time"
)

func (a *DBAdapter) GetRound(round_id string) (*model.Round, error) {
	players, err := a.GetRoundPlayers(round_id)
	if err != nil {
		return nil, err
	}

	comments, err := a.GetRoundComments(round_id)
	if err != nil {
		return nil, err
	}

	stmt := `SELECT * FROM Rounds WHERE id == ?;`

	rows, err := a.db.Query(stmt, round_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		round := &model.Round{}
		round.Players = players
		round.Comments = comments

		started_at := 0
		rows.Scan(&round.ID, &round.State, &started_at)
		round.StartedAt = time.Unix(int64(started_at), 0).Format(time.RFC3339)

		return round, nil
	}

	return nil, fmt.Errorf("round not found")
}

func (a *DBAdapter) ListRounds(n int, state model.RoundState) ([]*model.Round, error) {
	stmt := `
	SELECT * FROM Rounds
	WHERE state == ?
	ORDER BY started_at DESC
	LIMIT ?;
	`

	rows, err := a.db.Query(stmt, state, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rounds := make([]*model.Round, 0)

	for rows.Next() {
		round := &model.Round{}

		started_at := 0
		rows.Scan(&round.ID, &round.State, &started_at)
		round.StartedAt = time.Unix(int64(started_at), 0).Format(time.RFC3339)

		players, err := a.GetRoundPlayers(round.ID)
		if err != nil {
			return nil, err
		}

		comments, err := a.GetRoundComments(round.ID)
		if err != nil {
			return nil, err
		}

		round.Players = players
		round.Comments = comments

		rounds = append(rounds, round)
	}

	return rounds, nil
}

func (a *DBAdapter) InsertRound(info *model.RoundInfo) error {
	stmt := `
	INSERT INTO Rounds (id, state, started_at)
	VALUES (?, ?, ?);
	`

	started_at_t, err := time.Parse(time.RFC3339, info.StartedAt)
	if err != nil {
		return err
	}

	started_at := started_at_t.Unix()

	_, err = a.db.Exec(stmt, info.ID, info.State, started_at)

	return err
}

func (a *DBAdapter) UpdateRound(round_id string, state *model.RoundState) error {
	stmt := `
	UPDATE Rounds
	SET state = ?
	WHERE id == ?;
	`

	_, err := a.db.Exec(stmt, state, round_id)

	return err
}
