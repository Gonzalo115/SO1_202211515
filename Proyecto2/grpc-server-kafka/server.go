package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	pb "grpc-server/proto"

	"github.com/segmentio/kafka-go"

	"google.golang.org/grpc"
)

type server struct {
    pb.UnimplementedWeatherServiceServer
    writer *kafka.Writer
}

func (s *server) PostToKafka(ctx context.Context, req *pb.WeatherData) (*pb.WeatherResponse, error) {
	
	// Serializar WeatherData a JSON
	data, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error serializando WeatherData: %v", err)
		return nil, fmt.Errorf("error serializando datos: %v", err)
	}

	// Enviar mensaje a Kafka
	err = s.writer.WriteMessages(ctx,
		kafka.Message{
			Value: data,
		},
	)
	if err != nil {
		log.Printf("Error enviando mensaje a Kafka: %v", err)
		return nil, fmt.Errorf("error enviando mensaje a Kafka: %v", err)
	}
    log.Printf("Mensaje enviado a Kafka: %s", string(data))


	message := fmt.Sprintf("Datos Recibido correctamente del Country: %s", req.Country)

	return &pb.WeatherResponse{Message: message}, nil
}


func main() {
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "my-topic"
	}
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "my-cluster-kafka-bootstrap.kafka.svc.cluster.local:9092"
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		// Configuración de escritura
		BatchTimeout: 50 * time.Millisecond,
		WriteTimeout: 1 * time.Second,
		RequiredAcks: kafka.RequireOne,
		Logger:      kafka.LoggerFunc(log.Printf),
		ErrorLogger: kafka.LoggerFunc(log.Printf),
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpServer := grpc.NewServer()
	pb.RegisterWeatherServiceServer(grpServer, &server{writer: writer})
	log.Printf("Servidor gRPC Kafka corriendo en :50051")
	if err := grpServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	if err := writer.Close(); err != nil {
		log.Printf("Error cerrando writer de Kafka: %v", err)
	}

}

