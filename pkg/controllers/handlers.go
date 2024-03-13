package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/sdblg/meme/pkg/models"
	"github.com/sdblg/meme/pkg/services"
	"github.com/sdblg/meme/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
)

// Home displays the status of the api, as JSON.
func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Go up and running",
		Version: "1.0.0",
	}

	_ = utils.WriteJSON(w, http.StatusOK, payload)
}

// AllMemes returns a slice of all memes as JSON.
func (app *Application) AllMemes(w http.ResponseWriter, r *http.Request) {
	memes, err := app.DB.AllMemes()
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, memes)
}

// authenticate authenticates a user, and returns a JWT.
func (app *Application) authenticate(w http.ResponseWriter, r *http.Request) {
	// read json payload
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := utils.ReadJSON(w, r, &requestPayload)
	if err != nil {
		_ = utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate user against database
	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil {
		_ = utils.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// check password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		_ = utils.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// create a jwt user
	u := services.JwtUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	// generate tokens
	tokens, err := app.Auth.GenerateTokenPair(&u)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	refreshCookie := app.Auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	_ = utils.WriteJSON(w, http.StatusAccepted, tokens)
}

// refreshToken checks for a valid refresh cookie, and returns a JWT if it finds one.
func (app *Application) refreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == app.Auth.CookieName {
			claims := &services.Claims{}
			refreshToken := cookie.Value

			// parse the token to get the claims
			_, err := jwt.ParseWithClaims(
				refreshToken,
				claims,
				func(token *jwt.Token) (interface{}, error) {
					return []byte(app.JWTSecret), nil
				},
			)
			if err != nil {
				_ = utils.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			// get the user id from the token claims
			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				_ = utils.ErrorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			user, err := app.DB.GetUserByID(userID)
			if err != nil {
				_ = utils.ErrorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			u := services.JwtUser{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			}

			tokenPairs, err := app.Auth.GenerateTokenPair(&u)
			if err != nil {
				_ = utils.ErrorJSON(w, errors.New("error generating tokens"), http.StatusUnauthorized)
				return
			}

			http.SetCookie(w, app.Auth.GetRefreshCookie(tokenPairs.RefreshToken))

			_ = utils.WriteJSON(w, http.StatusOK, tokenPairs)

		}
	}
}

// logout logs the user out by sending an expired cookie to delete the refresh cookie.
func (app *Application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, app.Auth.GetExpiredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}

// GetMeme returns one meme, as JSON.
func (app *Application) GetMeme(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	memeID, err := strconv.Atoi(id)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	meme, err := app.DB.OneMeme(memeID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, meme)
}

// InsertMeme receives a JSON payload and tries to insert a meme into the database.
func (app *Application) InsertMeme(w http.ResponseWriter, r *http.Request) {
	var meme models.Meme

	err := utils.ReadJSON(w, r, &meme)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	meme.CreatedAt = time.Now()
	meme.UpdatedAt = time.Now()

	newID, err := app.DB.InsertMeme(meme)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("meme insterted with id: %d", newID),
	}

	_ = utils.WriteJSON(w, http.StatusAccepted, resp)
}

// UpdateMeme updates a meme in the database, based on a JSON payload.
func (app *Application) UpdateMeme(w http.ResponseWriter, r *http.Request) {
	var payload models.Meme

	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	meme, err := app.DB.OneMeme(payload.ID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	meme.Title = payload.Title
	meme.ReleaseDate = payload.ReleaseDate
	meme.Description = payload.Description
	meme.MPAARating = payload.MPAARating
	meme.RunTime = payload.RunTime
	meme.UpdatedAt = time.Now()

	err = app.DB.UpdateMeme(*meme)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{
		Error:   false,
		Message: "meme updated",
	}

	_ = utils.WriteJSON(w, http.StatusAccepted, resp)
}

// DeleteMeme deletes a meme from the database, by ID.
func (app *Application) DeleteMeme(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	err = app.DB.DeleteMeme(id)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{
		Error:   false,
		Message: "meme deleted",
	}

	_ = utils.WriteJSON(w, http.StatusAccepted, resp)
}

func EnableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://learn-code.ca")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().
				Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization")
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
}
