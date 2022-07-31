package main

import (
	"context"
	"fmt"
	"log"
	"loggerservice/data"
	"loggerservice/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer // ensure backwards compatibility
	Models                             data.Models
}

// WriteLog writes to Mongo via gRPC.
func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	// attempt to write entry to Mongo
	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{
			Result: err.Error(),
		}

		return res, err
	}

	// write back the response
	res := &logs.LogResponse{
		Result: "logged via gRPC",
	}

	return res, nil
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalln(err)
	}

	s := grpc.NewServer()
	logs.RegisterLogServiceServer(s, &LogServer{
		Models: app.Models,
	})

	log.Printf("gRPC server started on port %s", grpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
