package models

type User struct {
	Model
	Email     string
	AuthRead  bool
	AuthWrite bool
	AuthAdmin bool
	Name      string
	Notes     string
}

func UserByEmail(email string) User {
	u := User{}
	db.Where(&User{Email: email}).First(&u)
	return u
}
