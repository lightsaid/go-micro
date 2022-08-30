package main

import (
	"log"
	"net/http"
)

// NOTE:
// 用于发送电子邮件的 Golang 包。支持保持连接、TLS 和 SSL。易于批量 SMTP。
// go get github.com/xhit/go-simple-mail/v2

// golang中html邮件的内联样式
// go get github.com/vanng822/go-premailer

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