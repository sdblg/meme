package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application) Routes() http.Handler {
	// create a router mux
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(EnableCORS)

	mux.Get("/", app.Home)

	mux.Post("/authenticate", app.authenticate)
	mux.Get("/refresh", app.refreshToken)
	mux.Get("/logout", app.logout)

	mux.Get("/memes", app.AllMemes)
	mux.Get("/memes/{id}", app.GetMeme)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.Auth.AuthRequired)

		mux.Get("/memes/{id}", app.GetMeme)
		mux.Put("/memes", app.InsertMeme)
		mux.Patch("/memes/{id}", app.UpdateMeme)
		mux.Delete("/memes/{id}", app.DeleteMeme)
	})

	return mux
}
