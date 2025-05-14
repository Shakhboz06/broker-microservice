package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Mailer Mail
}

func main() {
	app := Config{
		Mailer: createMail(),
	}

	log.Println("Connecting to Mail Service ... on port 80")

	srv := &http.Server{
		Addr:    ":80",
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		// log.Printf("Server has caught error: %s", err)
		log.Panic(err)
	}
}

func createMail() Mail {

	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	
	mail := Mail{
		Domain: os.Getenv("MAIL_DOMAIN"),
		Host: os.Getenv("MAIL_HOST"),
		Port: port,
		Username: os.Getenv("MAIL_USERNAME"),
		Password: os.Getenv("MAIL_PASSWORD"),
		Encryption: os.Getenv("MAIL_ENCRYPTION"),
		FromAddress: os.Getenv("MAIL_FROM_ADDRESS"),
		FromName: os.Getenv("MAIL_FROM_NAME"),
	}
	
	return mail
}
