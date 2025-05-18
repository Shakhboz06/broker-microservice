package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) AuthUserTest(w http.ResponseWriter, r *http.Request) {

	var userPayload Payload
	if err := app.readJSON(w, r, &userPayload); err != nil {
		app.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	
	user, err := app.Repo.GetByEmail(ctx, userPayload.Email)

	if err != nil {
		app.ErrorResponse(w, errors.New("Invalid credentials"), http.StatusUnauthorized)
		return
	}

	valid, err := user.PasswordMatches(userPayload.Password)

	if err != nil || !valid {
		app.ErrorResponse(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.ErrorResponse(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("logged in user %s", user.Email),
		Data:    user,
	}
	

	app.WriteJSON(w, http.StatusAccepted, payload)
	

}
