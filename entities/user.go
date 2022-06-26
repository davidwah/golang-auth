package entities

import "time"

type User struct {
	Id        int
	Nama      string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}
