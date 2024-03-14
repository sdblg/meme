package models

import "time"

type Meme struct {
	ID        int       `json:"id"`
	Lan       string    `json:"lat"`
	Lon       string    `json:"lon"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
