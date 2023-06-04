package db

import (
	"grpc-mafia/round-tracker/graph/model"

	"github.com/google/uuid"
)

func (a *DBAdapter) GetRoundComments(round_id string) ([]*model.Comment, error) {
	stmt := `SELECT * FROM Comments WHERE round_id == ?;`

	rows, err := a.db.Query(stmt, round_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]*model.Comment, 0)

	for rows.Next() {
		comment := &model.Comment{}
		id := ""
		rid := ""

		rows.Scan(&rid, &id, &comment.From, &comment.Text)

		comments = append(comments, comment)
	}

	return comments, nil
}

func (a *DBAdapter) InsertRoundComment(round_id string, from string, text string) error {
	stmt := `
	INSERT INTO Comments (id, round_id, player_login, text)
	VALUES (?, ?, ?, ?);
	`

	id := uuid.New().String()

	_, err := a.db.Exec(stmt, id, round_id, from, text)

	return err
}
