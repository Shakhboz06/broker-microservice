package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	// "os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml")
	})

	fmt.Println("Starting front end service on port 8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Panic(err)
	}
}

//go:embed cmd/web/templates/*
var templateFS embed.FS


func render(w http.ResponseWriter, t string) {

	partials := []string{
		"cmd/web/templates/base.layout.gohtml",
		"cmd/web/templates/header.partial.gohtml",
		"cmd/web/templates/footer.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("cmd/web/templates/%s", t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	tmpl, err := template.ParseFS(templateFS, templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data struct{
		BrokerURL string
	}

	// data.BrokerURL = os.Getenv("BROKER_URL")
	data.BrokerURL = "http://localhost:8080"


	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
