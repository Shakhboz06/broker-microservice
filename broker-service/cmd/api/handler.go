package main

import (
	"broker-service/event"
	"broker-service/logs"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/rpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Action   string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}
type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type RpcPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Handler(w http.ResponseWriter, r *http.Request) {

	payload := jsonResponse{
		Error:   false,
		Message: "Success message from Broker",
	}

	err := app.WriteJSON(w, http.StatusOK, payload)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// w.Write(res)
}

func (app *Config) HandleSubmissions(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	if err := app.readJSON(w, r, &requestPayload); err != nil {
		app.ErrorResponse(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		// app.logItem(w, requestPayload.Log)
		// app.logEventwithRabbitMQ(w, requestPayload.Log)
		app.LogItemWithRPC(w, requestPayload.Log)
	case "mail":
		app.sendEmail(w, requestPayload.Mail)
	default:
		app.ErrorResponse(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// 1) Marshal the credentials into JSON
	

	jsonData, err := json.Marshal(a)
	if err != nil {
		app.ErrorResponse(w, fmt.Errorf("failed to marshal auth payload: %w", err), http.StatusInternalServerError)
		return
	}

	// 2) Build the HTTP request to the auth service
	req, err := http.NewRequest(
		http.MethodPost,
		"http://authentication-service/authentication",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		app.ErrorResponse(w, fmt.Errorf("failed to create auth request: %w", err), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// 3) Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		app.ErrorResponse(w, fmt.Errorf("auth service unreachable: %w", err), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// // 4) Handle HTTP status codes
	switch resp.StatusCode {
	case http.StatusOK, http.StatusAccepted:
		// good, continue
	case http.StatusUnauthorized:
		app.ErrorResponse(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		// app.ErrorResponse(w, err, http.StatusUnauthorized)
		return
	default:
		app.ErrorResponse(w, fmt.Errorf("auth service error: status %d", resp.StatusCode), http.StatusBadGateway)
		return
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusAccepted:
	case http.StatusUnauthorized:
		app.ErrorResponse(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	default:
		app.ErrorResponse(w,
			fmt.Errorf("auth service error: status %d", resp.StatusCode),
			http.StatusBadGateway,
		)
		return
	}

	var upstream jsonResponse
	if err := json.NewDecoder(resp.Body).Decode(&upstream); err != nil {
		app.ErrorResponse(w, fmt.Errorf("failed to decode auth response: %w", err), http.StatusBadGateway)
		return
	}

	if upstream.Error {
		app.ErrorResponse(w, errors.New(upstream.Message), http.StatusUnauthorized)
		return
	}

	out := jsonResponse{
		Error:   false,
		Message: "Authenticated!",
		Data:    upstream.Data,
	}
	if err := app.WriteJSON(w, http.StatusOK, out); err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

// func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {

// 	jsonData, _ := json.Marshal(entry)

// 	logServiceURL := "http://logger-service/log"

// 	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		app.ErrorResponse(w, err)
// 		return
// 	}

// 	request.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}

// 	response, err := client.Do(request)
// 	if err != nil {
// 		app.ErrorResponse(w, err)
// 		return
// 	}
// 	defer response.Body.Close()

// 	if response.StatusCode != http.StatusAccepted {
// 		app.ErrorResponse(w, err)
// 		return
// 	}

// 	var payload jsonResponse
// 	payload.Error = false
// 	payload.Message = "logged"

// 	app.WriteJSON(w, http.StatusAccepted, payload)

// }

func (app *Config) sendEmail(w http.ResponseWriter, mail MailPayload) {

	jsonData, _ := json.Marshal(mail)

	mailServiceURL := "http://mailer-service/send"
	req, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.ErrorResponse(w, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(req)

	if err != nil {
		app.ErrorResponse(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.ErrorResponse(w, errors.New("mail service error"))
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "sent to" + mail.To

	app.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logEventwithRabbitMQ(w http.ResponseWriter, log LogPayload) {

	err := app.pushToQueue(log.Name, log.Data)
	if err != nil {
		app.ErrorResponse(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Logged Successfully with RabbitMQ"

	app.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.RabbitMQ)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	jsonData, _ := json.Marshal(&payload)

	err = emitter.Push(string(jsonData), "log.INFO")
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) LogItemWithRPC(w http.ResponseWriter, log LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		app.ErrorResponse(w, err)
		return
	}

	rpcPayload := RpcPayload{
		Name: log.Name,
		Data: log.Data,
	}

	var result string

	err = client.Call("RpcServer.LogInfo", rpcPayload, &result)
	if err != nil {
		app.ErrorResponse(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: result,
	}

	app.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logWithgRPC(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	if err := app.readJSON(w, r, &requestPayload); err != nil {
		app.ErrorResponse(w, err)
		return
	}

	// conn, err := grpc.Dial("logger-service:5001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	conn, err := grpc.NewClient("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		app.ErrorResponse(w, err)
		return
	}

	defer conn.Close()

	client := logs.NewLogServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = client.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})

	if err != nil {
		app.ErrorResponse(w, err)
		return
	}

	var payload jsonResponse

	payload.Error = false
	payload.Message = "Logged successfully with gRPC"

	app.WriteJSON(w, http.StatusAccepted, payload)

}
