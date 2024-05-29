package models

import "time"

type URL struct {
	Id        int64     `json:"id"`
	Url       string    `json:"url"`
	Alias     string    `json:"alias"`
	CreatedAt time.Time `json:"createdAt"`
}

type User struct {
	Id        int64     `json:"id"`
	Username  string    `json:"username" validate:"required,max=20"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=5"`
	CreatedAt time.Time `json:"createdAt"`
}

type LoginData struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
