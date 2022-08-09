package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"lightsaid.com/go-micro/auth-service/data"
)

func (app *Config) Insert(w http.ResponseWriter, r *http.Request) {
	var req struct{
		Email string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}

	err := app.readJSON(w, r, &req)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user := data.User{
		Email: req.Email,
		Password: req.Password,
		Username: req.Username,
	}

	_, err = user.Insert(user)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error: false,
		Message: fmt.Sprintf("创建用户成功"),
		Data: user,
	}
	
	app.writeJSON(w, http.StatusOK, payload)

}

// Auth 认证服务
func (app *Config) Auth(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &req)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(req.Email)
	if err != nil {
		app.errorJSON(w, errors.New("无效的Email"), http.StatusBadRequest)
		return
	}

	_, err = user.PasswordMatches(req.Password)
	if err != nil {
		app.errorJSON(w, errors.New("认证无效，密码不对"))
		return
	}

	// log auth
	err = app.LogRequest("auth", fmt.Sprintf("%s logged in ", user.Email))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error: false,
		Message: fmt.Sprintf("登录用户： %s", user.Email),
		Data: user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) LogRequest(name, data string) error{
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "","\t")
	url := "http://logger-service/log"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)

	return nil
}  
