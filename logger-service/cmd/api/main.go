package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webport  = ":80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo_db:27017"
	gRpcPort = "50001"
)

type Config struct {
	DB     *mongo.Client
	Models data.Models
}

func main() {

	app := Config{}
	var mongoClient *mongo.Client
	var err error
	for attempt, backoff := 1, time.Second; attempt <= 5; attempt, backoff = attempt+1, backoff*2 {
		mongoClient, err = app.connectToMongoDB()
		if err == nil {
			break
		}
		log.Printf("MongoDB connect attempt %d failed: %v; retrying in %s\n", attempt, err, backoff)
		time.Sleep(backoff)
	}
	if err != nil {
		log.Fatalf("Could not connect to MongoDB after retries: %v", err)
	}
	
	log.Println("Connected to MongoDB")
	
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Fatalf("Error disconnecting MongoDB client: %v", err)
		}
	}()

	app.DB = mongoClient
	app.Models = *data.New(app.DB)

	// Register RPC Server
	rpcSrv := &RpcServer{client: app.DB}
	if err := rpc.Register(rpcSrv); err != nil {
		log.Fatalf("RPC Register failed: %v", err)
	}
	go app.rpcListen()

	go app.gRPCListen()
	
	log.Println("Starting the Log server on Port", webport)

	srv := &http.Server{
		Addr:    webport,
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panicf("Could not connect to log server %v", err)
	}
}

func (app *Config) connectToMongoDB() (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI(mongoURL).SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	}).SetMaxPoolSize(50).
		SetMinPoolSize(5).
		SetMaxConnIdleTime(30 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	conn, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("Error Connecting to MongoDB", err)
		return nil, err
	}

	//this is just checkinf whether it is connecting with Ping
	if err := conn.Ping(ctx, nil); err != nil {
		conn.Disconnect(ctx)
		return nil, fmt.Errorf("mongo Pinging: %w", err)
	}

	return conn, nil
}

func (app *Config) rpcListen() {
	log.Println("Starting RPC Server on port ..", rpcPort)

	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		log.Panic(err)
	}

	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}

		go rpc.ServeConn(rpcConn)
	}
}
