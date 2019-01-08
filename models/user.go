package models

type User struct {
	ID       int    `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password,omitempty"`
	Admin    bool   `db:"admin" json:"admin"`
	Quota    *int   `db:"quota" json:"quota,omitempty"`
}
