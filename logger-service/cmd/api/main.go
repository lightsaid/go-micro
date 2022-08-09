package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"lightsaid.com/go-micro/logger-service/data"
)


const (
	webPort = ":80"
	// rpcPort = ":5001"
	// gRpcPort = ":50001"
	mongoURL = "mongodb://mongo:27017"
)

var client *mongo.Client

type Config struct{
	Models data.Models
}

func main(){
	// 链接 mongo
	momngoClient, err := openMongo()
	if err != nil {
		log.Panic(err)
	}
	client = momngoClient

	// 创建上下关闭mongo链接
	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
	defer cancel()

	// 关闭链接
	defer func(){
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := &Config{
		Models: data.New(client),
	}

	app.serve()
}

func (app *Config) serve() {
	srv := http.Server{
		Addr: webPort,
		Handler: app.routes(),
	}

	log.Println("Starting server on port ", webPort)
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func openMongo() (*mongo.Client, error){
	// 创建链接参数
	opts := options.Client().ApplyURI(mongoURL)
	opts.SetAuth(options.Credential{
		Username: "admin",
		Password: "abc123",
	})

	// 链接
	c, err := mongo.Connect(context.TODO(), opts)
	if err != nil{
		log.Println("Error connecting: ", err)
		return nil, err
	}

	return c, nil
}


