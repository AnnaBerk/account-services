package entity

import "time"

type User struct {
	Id        int       `db:"id"`
	Username  string    `db:"name"`
	Password  string    `db:"password"`
	AccountId int       `account_id:"id"`
	CreatedAt time.Time `db:"created"`
}
