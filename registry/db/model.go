package db

type User struct {
	Login          string
	AvatarFilename string
	Gender         string
	Mail           string
}

func (u *User) Merge(o *User) {
	if u.Login != o.Login {
		panic("merge must be done on the same user")
	}

	if len(u.AvatarFilename) == 0 {
		u.AvatarFilename = o.AvatarFilename
	}

	if len(u.Gender) == 0 {
		u.Gender = o.Gender
	}

	if len(u.Mail) == 0 {
		u.Mail = o.Mail
	}
}
