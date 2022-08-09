package main

import (
	"log"
	"net/http"
)

const webPort = ":80"
type Config struct{}

func main(){
	app := Config{}

	log.Println("Starting mail service on ", webPort)

	srv := &http.Server{
		Addr: webPort,
		Handler: app.routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}

}