package main

import (
	"context"
	"fmt"
	"log"
	"loggerservice/data"
	"time"
)

type RPCServer struct{}

// RPCPayload is the payload we expect to receive from RPC.
type RPCPayload struct {
	Name string
	Data string
}

// LogInfo logs the given payload to Mongo.
func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println(err)
		return err
	}

	*resp = fmt.Sprintf("Processed payload via RPC: %s", payload.Name)

	return nil
}
