package main

import (
	"net/http"

	"lightsaid.com/go-micro/logger-service/data"
)

type JSONPayload struct{
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request){
	var req JSONPayload
	_ = app.readJSON(w, r, &req)

	l := data.LogEntry{
		Name: req.Name,
		Data: req.Data,
	}

	err := app.Models.LogEntry.Insert(l)
	if err != nil {
		app.errorJSON(w, err)
		return 
	}

	resp := jsonResponse {
		Error: false,
		Message: "Write Log Success.",
	}
	app.writeJSON(w, http.StatusAccepted, resp)
}