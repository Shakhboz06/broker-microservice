package main

import (
	"context"
	"log"
	"logger-service/data"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type RpcServer struct {
	client *mongo.Client
}

type RpcPayload struct {
	Name string
	Data string
}

func (r *RpcServer) LogInfo(pay RpcPayload, resp *string) error {
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	collection := data.NewLogStore(r.client, "appdb", "logs")
	
	err := collection.Insert(ctx, &data.LogEntry{
			Name: pay.Name,
			Data: pay.Data,
			CreatedAt: time.Now(),
	})

	if err != nil {
		log.Println("Error writing to mongo database", err)
		return err
	}

	*resp = "Processed payload with RPC" + pay.Name

	return nil
}
