package models

import "time"

type Token struct {
	Token   string    `db:"token"`
	Expires time.Time `db:"expires"`
	Scope   string    `db:"scope"`
	UserID  int       `db:"user_id"`
}
