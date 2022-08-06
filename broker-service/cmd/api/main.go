package main

import (
	"log"
	"net/http"
)

const webPort = ":8080"

type Config struct{}

func main(){
	app := Config{}

	log.Println("Starting broker service on port ", webPort)

	srv := &http.Server{
		Addr: webPort,
		Handler: app.routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err.Error())
	}
}
