package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "grpc-server/proto"

	"google.golang.org/grpc"
)

type server struct {
    pb.UnimplementedWeatherServiceServer
}

func (s *server) PostToRabbitMQ(ctx context.Context, req *pb.WeatherData) (*pb.WeatherResponse, error) {
	
    //Guardar en kafka

	message := fmt.Sprintf("Datos Recibido correctamente del Country: %s", req.Country)

	return &pb.WeatherResponse{Message: message}, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50052")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    grpServer := grpc.NewServer()
    pb.RegisterWeatherServiceServer(grpServer, &server{})
    if err := grpServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}