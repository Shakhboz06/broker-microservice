package main

import (
	"logger-service/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)

	ctx := r.Context()
	
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}
	
	err := app.Models.Logs.Insert(ctx, &event)
	if err != nil {
		app.ErrorResponse(w, err)
		return
	}

	resp := jsonResponse{
		Error: false,
		Message: "logged",
	}

	app.WriteJSON(w, http.StatusAccepted, resp)
}

