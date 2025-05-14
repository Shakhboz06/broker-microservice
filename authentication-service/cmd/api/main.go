package main

import (
	"authentication-service/data"
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	DB *sql.DB
	Users *data.UserStore
}

func main(){
	log.Println("Starting Authentication Server")

	conn := connectToDB()

	// if conn == nil {
	// 	log.Panic("Can't connect to Postgres!")
	// }

	app := &Config{
 		DB: conn,
		Users: data.NewUserStore(conn),
	}


	srv := &http.Server{
		Addr: ":80",
		Handler: app.routes(),

	}


	log.Printf("Listening on %s\n", srv.Addr)
	err := srv.ListenAndServe()
	if err != nil{
		log.Panic(err)
	}

}



func openDB(addr string) (*sql.DB, error){

	db, err := sql.Open("pgx", addr)

	if err != nil{
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()


	err = db.PingContext(ctx)

	if err != nil{
		return nil, err
	}

	return db, nil
}


func connectToDB() *sql.DB{

 	dsn := os.Getenv("DSN")

	for i := 1; i <= 10; i++{

		conn, err := openDB(dsn)
		
		if err == nil {
			log.Println("Connected to Postgres!")
			return conn
		}

		log.Printf("Postgres not ready (%d/10): %v", i, err)
		time.Sleep(2 * time.Second)

	}

	log.Panic("Could not connect to Postgres after 10 attempts")

	return nil
}