package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// RequestPayload 代理服务器，请求类型结构体, 根据Action分发请求处理
type RequestPayload struct{
	Action string `json:"action"`
	Auth AuthPayload `json:"auth,omitempty"`
	User UserPayload `json:"user"`
}

type UserPayload struct {
	Email string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type AuthPayload struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error: false,
		Message: "Hello broker service",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

}


func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var req RequestPayload

	err := app.readJSON(w, r, &req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	 switch req.Action {
	 case "auth":
		app.authenticate(w, req.Auth)
	 case "add_user":
		app.insert(w, req.User)
	 default:
		app.errorJSON(w, errors.New("unknow action"))
	 }

}

func (app *Config) insert(w http.ResponseWriter, data UserPayload){
	arg, _ := json.Marshal(data)

	// 创建一个请求，向 docker compose 中 auth-service 服务发请求
	url := "http://auth-service:8080/insert"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(arg))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// 实例化一个HTTP客户端
	client := &http.Client{}

	// HTTP 客户端发送请求
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	fmt.Printf("代理请求: %s 发送成功。请求参数：%s, 响应状态码：%v\n", url, arg, response.StatusCode)

	if response.StatusCode != http.StatusOK{
		log.Println("创建用户失败，状态码：",response.StatusCode)
		app.errorJSON(w, errors.New("创建用户失败"))
		return
	}

	// 检查响应数据是否包含错误
	var responseData jsonResponse
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	fmt.Printf("代理请求: %s, 返回数据：%v", url, responseData)

	if responseData.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	// 成功
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Auth Success!"
	payload.Data = responseData.Data

	app.writeJSON(w, http.StatusAccepted, payload)

}


func (app *Config) authenticate(w http.ResponseWriter, data AuthPayload) {
	jsonData, _ := json.Marshal(data)

	// 创建一个请求，向 docker compose 中的 auth-service 服务发请求
	url := "http://auth-service:8080/auth"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// 实例化一个HTTP客户端
	client := &http.Client{}

	// HTTP 客户端发送请求
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()
	fmt.Printf("代理请求: %s 发送成功。请求参数：%s, 响应状态码：%v\n", url, jsonData, response.StatusCode)

 	if response.StatusCode != http.StatusAccepted{
		log.Println("认证不成功，请求状态码：",response.StatusCode)
		app.errorJSON(w, errors.New("认证不成功"))
		return
	}

	// 成功
	var responseData jsonResponse
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	fmt.Printf("代理请求: %s, 返回数据：%v", url, responseData)

	if responseData.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Auth Success!"
	payload.Data = responseData.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}