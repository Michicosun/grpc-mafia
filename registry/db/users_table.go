package db

func (a *DBAdapter) PrefetchUser(login string) error {
	stmt := `INSERT OR IGNORE INTO Users (login) VALUES (?);`
	_, err := a.db.Exec(stmt, login)
	return err
}

func (a *DBAdapter) GetUserOrCreateDefault(login string) (*User, error) {
	a.PrefetchUser(login)

	stmt := `SELECT * FROM Users WHERE login == ?;`
	rows, err := a.db.Query(stmt, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := User{}
		rows.Scan(&user.Login, &user.AvatarFilename, &user.Gender, &user.Mail)
		return &user, nil
	}

	panic("unreachable")
}

func (a *DBAdapter) UpdateUser(new User) (*User, error) {
	cur, err := a.GetUserOrCreateDefault(new.Login)
	if err != nil {
		return nil, err
	}

	new.Merge(cur)

	stmt := `
	UPDATE Users
	SET avatar_filename = ?,
		gender = ?,
		mail = ?
	WHERE login == ?;
	`

	_, err = a.db.Exec(stmt, new.AvatarFilename, new.Gender, new.Mail, new.Login)

	return &new, err
}

func (a *DBAdapter) DeleteUser(login string) error {
	stmt := `DELETE FROM Users WHERE login == ?;`
	_, err := a.db.Exec(stmt, login)
	return err
}
