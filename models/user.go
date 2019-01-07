package models

type User struct {
	ID       int    `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Admin    bool   `db:"admin"`
	Quota    int    `db:"quota"`
}
