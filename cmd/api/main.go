package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sdblg/meme/pkg/controllers"
	"github.com/sdblg/meme/pkg/repository/dbrepo"
	"github.com/sdblg/meme/pkg/services"
)

const port = 8080

func main() {
	// set Application config
	var app controllers.Application

	// read from command line
	flag.StringVar(
		&app.DSN,
		"dsn",
		"host=localhost port=5432 user=esusu password=esusu dbname=postgres sslmode=disable timezone=UTC connect_timeout=5",
		"Postgres connection string",
	)
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "signing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "esusu.com", "signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "esusu.com", "signing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "cookie domain")
	flag.StringVar(&app.Domain, "domain", "esusu.com", "domain")
	flag.Parse()

	// connect to the database
	conn, err := app.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()

	app.Auth = services.Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "__Host-refresh_token",
		CookieDomain:  app.CookieDomain,
	}

	log.Println("Starting Application on port", port)

	// start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.Routes())
	if err != nil {
		log.Fatal(err)
	}
}
