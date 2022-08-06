package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"lightsaid.com/go-micro/auth-service/data"

	_ "github.com/lib/pq"
)

const webPort = ":8080"

type Config struct {
	DB *sql.DB
	Models data.Models
}

func main(){

	dsn := os.Getenv("DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	app := Config{
		DB: db,
		Models: data.New(db),
	}

	srv := &http.Server{
		Addr: webPort,
		Handler: app.routes(),
	}

	log.Println("Starting autu server on port ", webPort)

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

