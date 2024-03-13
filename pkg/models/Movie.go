package models

import "time"

type Meme struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	ReleaseDate time.Time `json:"release_date"`
	RunTime     int       `json:"runtime"`
	MPAARating  string    `json:"mpaa_rating"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Genres      []*Genre  `json:"categories,omitempty"`
	GenresArray []int     `json:"genres_array,omitempty"`
}

type Genre struct {
	ID        int       `json:"id"`
	Genre     string    `json:"genre"`
	Checked   bool      `json:"checked"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
