package models

import "time"

type URL struct {
	Id        int64     `json:"id"`
	Url       string    `json:"url"`
	Alias     string    `json:"alias"`
	CreatedAt time.Time `json:"createdAt"`
}
