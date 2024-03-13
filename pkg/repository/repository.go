package repository

import (
	"database/sql"

	"github.com/sdblg/meme/pkg/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB

	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)

	AllMemes() ([]*models.Meme, error)
	OneMeme(id int) (*models.Meme, error)

	InsertMeme(meme models.Meme) (int, error)
	UpdateMeme(meme models.Meme) error
	DeleteMeme(id int) error
}
