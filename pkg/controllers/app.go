package controllers

import (
	"database/sql"
	"log"

	"github.com/sdblg/meme/pkg/repository"
	"github.com/sdblg/meme/pkg/services"
)

type Application struct {
	DSN          string
	Domain       string
	DB           repository.DatabaseRepo
	Auth         services.Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
}

func (app *Application) ConnectToDB() (*sql.DB, error) {
	connection, err := services.OpenDB(app.DSN)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Postgres!")
	return connection, nil
}
