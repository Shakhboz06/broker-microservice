package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error{
	 
	maxBytes := 10487576 //1 mb

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes)) 
	// http.MaxBytesReader is a function that  is intended for limiting the size of incoming request bodies.
	// MaxBytesReader prevents clients from accidentally or maliciously sending a large request and wasting server resources.


	dec := json.NewDecoder(r.Body)

	err := dec.Decode(data)
	if err != nil{
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must be a single JSON value")

	}


	return nil

}



func (app *Config) WriteJSON( w http.ResponseWriter, status int, data any, headers ...http.Header) error{
	

	output, err := json.Marshal(data)

	if err != nil{
		return err
	}

	if len(headers) > 0{
		for key, value := range headers[0]{
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err = w.Write(output); err != nil{
		return err
	}

	return nil

}



func (app *Config) ErrorResponse(w http.ResponseWriter, err error, status ...int) error{
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0] 
	}


	payload := &jsonResponse{
		Error: true,
		Message: err.Error(),
	}


	return app.WriteJSON(w, statusCode, payload)
}