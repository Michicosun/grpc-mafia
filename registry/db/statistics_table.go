package db

import (
	"fmt"
	"time"
)

type GameStatistics struct {
	SessionsCnt string `json:"sessions_cnt"`
	WinCnt      uint64 `json:"win_cnt"`
	LoseCnt     uint64 `json:"lose_cnt"`
	TotalTime   string `json:"total_time"`
}

func (a *DBAdapter) PrefetchStatistics(login string) error {
	stmt := `
	INSERT OR IGNORE INTO Statistics (
		login,
		sessions_cnt,
		win_cnt,
		lose_cnt,
		total_time_sec
	) VALUES (?, 0, 0, 0, 0);`

	_, err := a.db.Exec(stmt, login)
	return err
}

func (a *DBAdapter) AddNewRound(login string, win uint8, round_time time.Duration) error {
	if err := a.PrefetchStatistics(login); err != nil {
		return err
	}

	stmt := `
	UPDATE Statistics
	SET sessions_cnt = sessions_cnt + ?,
		win_cnt = win_cnt + ?,
		lose_cnt = lose_cnt + ?,
		total_time_sec = total_time_sec + ?
	WHERE login == ?;
	`
	_, err := a.db.Exec(stmt, 1, win, 1-win, uint64(round_time.Seconds()), login)

	return err
}

func (a *DBAdapter) GetStatistics(login string) (*GameStatistics, error) {
	stmt := `SELECT * FROM Statistics WHERE login == ?;`

	rows, err := a.db.Query(stmt, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		stat := GameStatistics{}
		sec_cnt := int64(0)

		rows.Scan(&login, &stat.SessionsCnt, &stat.WinCnt, &stat.LoseCnt, &sec_cnt)
		total_duration := time.Duration(int64(time.Second) * sec_cnt)

		hours := int64(total_duration.Hours())
		minutes := int64(total_duration.Minutes()) % 60
		seconds := int64(total_duration.Seconds()) % 60

		stat.TotalTime = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

		return &stat, nil
	}

	return nil, fmt.Errorf("statistics for user: %s not found", login)
}
