package main

import (
	"bytes" 
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RoundTripFunc func (req *http.Request) *http.Response


func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error){
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc)*http.Client{
	return &http.Client{
		Transport: fn,
	}
}
func Test_Authentication(t *testing.T) {

	jsonToReturn := `
		{
			"error": "false",
			"message": "some message"
		}	
	`
	client := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusAccepted,
			Body: io.NopCloser(bytes.NewBufferString(jsonToReturn)),
			Header: make(http.Header),
		}
	})

	testApp.Client = client 
	postBody := map[string]string{
		"email": "me@here.com",
		"password": "verysecret",
	}

	body, _ := json.Marshal(postBody)

	req, _ :=  http.NewRequest("POST", "/authentication", bytes.NewReader(body))
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testApp.AuthUserTest)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted{
		t.Errorf("expected http.Accepted but got %d", rr.Code)
	}
}