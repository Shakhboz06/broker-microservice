package main

import (
	"authentication-service/data"
	"os"
	"testing"
)

var testApp Config
func TestMain(m *testing.M){
	repo := data.NewPostGresTest(nil)
	testApp.Repo = repo
	os.Exit(m.Run())
}