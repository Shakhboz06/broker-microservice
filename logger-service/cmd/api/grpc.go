package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"logger-service/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (log *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest)(*logs.LogResponse, error){
	input := req.GetLogEntry()

	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := log.Models.Logs.Insert(ctx, &logEntry)
	if err != nil{
		res := &logs.LogResponse{
			Result: "Failed",
		}
		return res, err
	}

	res := &logs.LogResponse{
		Result: "logged!",
	}
	return res, nil
}



func(app *Config) gRPCListen(){
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil{
		log.Fatalf("Failed to listen to gRPC %v", err)
	}
	
	serv := grpc.NewServer()
	
	logs.RegisterLogServiceServer(serv, &LogServer{Models: app.Models})
	
	log.Printf("gRPC server starting on port ... %s", gRpcPort)
	
	if err := serv.Serve(listen); err != nil{
		log.Fatalf("Failed to listen to gRPC %v", err)
	}
	
}