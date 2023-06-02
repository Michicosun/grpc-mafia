package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DBAdapter struct {
	db *sql.DB
}

func (a *DBAdapter) InitTables() {
	users_sql_stmt := `
	CREATE TABLE IF NOT EXISTS Users (
		login TEXT NOT NULL PRIMARY KEY,
		avatar_filename TEXT,
		gender TEXT,
		mail TEXT
	);
	`
	_, err := a.db.Exec(users_sql_stmt)
	if err != nil {
		panic(err)
	}

	statistics_sql_stmt := `
	CREATE TABLE IF NOT EXISTS Statistics (
		login TEXT NOT NULL PRIMARY KEY,
		sessions_cnt INTEGER NOT NULL,
		win_cnt INTEGER NOT NULL,
		lose_cnt INTEGER NOT NULL,
		total_time_sec INTEGER NOT NULL
	);
	`
	_, err = a.db.Exec(statistics_sql_stmt)
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
