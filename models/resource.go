package models

import "time"

type Resource struct {
	ID        int       `db:"id" json:"-"`
	Key       string    `db:"key" json:"key"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UserID    int       `db:"user_id" json:"-"`
}
