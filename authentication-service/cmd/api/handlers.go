package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"


	"net/http"
)

type Payload struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

type Entry struct{
	Name string `json:"name"`
	Data string `json:"data,omitempty"`
}

func (app *Config) AuthUser(w http.ResponseWriter, r *http.Request) {

	var userPayload Payload
	if err := app.readJSON(w, r, &userPayload);

	err != nil {
		app.ErrorResponse(w,err, http.StatusBadRequest)
		return
	}
	
	
	ctx := r.Context()
	user, err := app.Users.GetByEmail(ctx, userPayload.Email)

	
	if err != nil {
		app.ErrorResponse(w, errors.New("Invalid credentials"), http.StatusUnauthorized)
		return
	}
	
	
	valid, err := user.PasswordMatches(userPayload.Password)

	if err != nil || !valid{
		app.ErrorResponse(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}
	
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil{
		app.ErrorResponse(w, err)
		return
	}
	
	payload := jsonResponse{
		Error: false,
		Message: fmt.Sprintf("logged in user %s", user.Email),
		Data: user,
	}

	app.WriteJSON(w, http.StatusOK, payload)


}

func(app *Config) logRequest(name, data string)error{

	entry := Entry{
		Name: name,
		Data: data,
	}

	jsonData, _ := json.Marshal(entry)
	logServiceUrl := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	
	

	_, err = app.Client.Do(request)
	if err != nil {
		return err
	}

	return nil
	
}