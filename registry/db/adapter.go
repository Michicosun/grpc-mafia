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
		login string not null primary key,
		avatar_filename string,
		gender string,
		mail string
	);
	`
	_, err := a.db.Exec(users_sql_stmt)
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
