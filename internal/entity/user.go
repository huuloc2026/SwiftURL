package entity

import "time"

type User struct {
	ID        int64     `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"` // hashed
	CreatedAt time.Time `db:"created_at"`
}
