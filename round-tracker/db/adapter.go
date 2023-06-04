package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DBAdapter struct {
	db *sql.DB
}

func (a *DBAdapter) InitTables() {
	rounds_sql_stmt := `
	CREATE TABLE IF NOT EXISTS Rounds (
		id TEXT NOT NULL PRIMARY KEY,
		state TEXT NOT NULL,
		started_at INTEGER NOT NULL
	);
	`
	_, err := a.db.Exec(rounds_sql_stmt)
	if err != nil {
		panic(err)
	}

	players_sql_stmt := `
	CREATE TABLE IF NOT EXISTS Players (
		round_id TEXT NOT NULL,
		login TEXT NOT NULL,
		role TEXT NOT NULL,
		alive INTEGER NOT NULL,
		PRIMARY KEY (round_id, login)
	);
	`
	_, err = a.db.Exec(players_sql_stmt)
	if err != nil {
		panic(err)
	}

	comments_sql_stmt := `
	CREATE TABLE IF NOT EXISTS Comments (
		id TEXT NOT NULL PRIMARY KEY,
		round_id TEXT NOT NULL,
		player_login TEXT NOT NULL,
		text TEXT NOT NULL
	);
	`
	_, err = a.db.Exec(comments_sql_stmt)
	if err != nil {
		panic(err)
	}
}

func NewDBAdapter(file string) *DBAdapter {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		panic(err) // can't work without db
	}

	return &DBAdapter{
		db: db,
	}
}
