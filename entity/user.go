package entity

import "time"

type User struct {
	Id        int       `json:"id"`
	Username  string    `json:"user_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
