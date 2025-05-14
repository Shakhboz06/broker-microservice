package main

import (
	"log"
	"net/http"
)

type mailMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"` 
}

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {

	var requestPayload mailMessage

	if err := app.readJSON(w, r, &requestPayload); err != nil{
		app.ErrorResponse(w, err)
		return
	}

	msg := Message{
		From: requestPayload.From,
		To: requestPayload.To,
		Subject: requestPayload.Subject,
		Data: requestPayload.Message,
	}

	err := app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		log.Println(err)
		app.ErrorResponse(w, err)
	}

	payload := jsonResponse{
		Error: false,
		Message: "Sent to " + requestPayload.To,
	}

	app.WriteJSON(w, http.StatusAccepted, payload) 

}
