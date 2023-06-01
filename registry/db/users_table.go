package db

import "fmt"

func (a *DBAdapter) prefetchUser(login string) error {
	stmt := `INSERT OR IGNORE INTO Users (login) VALUES (?);`
	_, err := a.db.Exec(stmt, login)
	return err
}

func (a *DBAdapter) GetUser(login string) (*User, error) {
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

	return nil, fmt.Errorf("user: %s not found", login)
}

func (a *DBAdapter) GetAllUsers() ([]User, error) {
	stmt := `SELECT * FROM Users;`

	rows, err := a.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		user := User{}
		rows.Scan(&user.Login, &user.AvatarFilename, &user.Gender, &user.Mail)
		users = append(users, user)
	}

	return users, nil
}

func (a *DBAdapter) UpdateUser(new User) (*User, error) {
	if err := a.prefetchUser(new.Login); err != nil {
		return nil, err
	}

	cur, err := a.GetUser(new.Login)
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
