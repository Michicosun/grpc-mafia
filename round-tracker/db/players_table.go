package db

import (
	"grpc-mafia/round-tracker/graph/model"
)

func (a *DBAdapter) GetRoundPlayers(round_id string) ([]*model.Player, error) {
	stmt := `SELECT * FROM Players WHERE round_id == ?;`

	rows, err := a.db.Query(stmt, round_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make([]*model.Player, 0)

	for rows.Next() {
		player := &model.Player{}
		rid := ""
		alive := int(0)

		rows.Scan(&rid, &player.Login, &player.Role, &alive)
		player.Alive = (alive == 1)

		players = append(players, player)
	}

	return players, nil
}

func (a *DBAdapter) InsertRoundPlayer(round_id string, info *model.PlayerInfo) error {
	stmt := `
	INSERT INTO Players (round_id, login, role, alive)
	VALUES (?, ?, ?, ?);
	`

	alive := 1

	_, err := a.db.Exec(stmt, round_id, info.Login, info.Role, alive)

	return err
}

func (a *DBAdapter) UpdateRoundPlayer(round_id string, status *model.PlayerStatus) error {
	stmt := `
	UPDATE Players
	SET alive = ?
	WHERE round_id == ? AND login == ?;
	`

	alive := 0
	if status.Alive {
		alive = 1
	}

	_, err := a.db.Exec(stmt, alive, round_id, status.Login)

	return err
}
