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

	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type server struct {
    pb.UnimplementedWeatherServiceServer
    rabbitConn *amqp091.Connection
	rabbitCh   *amqp091.Channel
	queueName  string
}

func (s *server) PostToRabbitMQ(ctx context.Context, req *pb.WeatherData) (*pb.WeatherResponse, error) {
    // Crear un contexto con timeout de 10 segundos
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    // Serializar WeatherData a JSON
    data, err := json.Marshal(req)
    if err != nil {
        log.Printf("Error serializando WeatherData: %v", err)
        return nil, fmt.Errorf("error serializando datos: %v", err)
    }

    // Publicar mensaje a RabbitMQ
    err = s.rabbitCh.PublishWithContext(ctx,
        "",
        s.queueName,
        false,
        false,
        amqp091.Publishing{
            ContentType: "application/json",
            Body:        data,
        },
    )
    if err != nil {
        log.Printf("Error publicando mensaje a RabbitMQ: %v", err)
        return nil, fmt.Errorf("error publicando mensaje a RabbitMQ: %v", err)
    }
    log.Printf("Mensaje enviado a RabbitMQ: %s", string(data))

	message := fmt.Sprintf("Datos Recibido correctamente del Country: %s", req.Country)

	return &pb.WeatherResponse{Message: message}, nil
}
func main() {
	// Configuración de RabbitMQ
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@rabbitmq.rabbitmq.svc.cluster.local:5672/"
	}
	queueName := os.Getenv("RABBITMQ_QUEUE")
	if queueName == "" {
		queueName = "weather-queue"
	}

	var rabbitConn *amqp091.Connection
	var rabbitCh *amqp091.Channel
	for i := 0; i < 5; i++ {
		conn, err := amqp091.Dial(rabbitMQURL)
		if err == nil {
			rabbitConn = conn
			break
		}
		log.Printf("Error conectando a RabbitMQ (intento %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if rabbitConn == nil {
		log.Fatalf("No se pudo conectar a RabbitMQ después de varios intentos")
	}
	defer rabbitConn.Close()

	// Crear canal
	rabbitCh, err := rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Error creando canal de RabbitMQ: %v", err)
	}
	defer rabbitCh.Close()

	// Declarar cola
	_, err = rabbitCh.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error declarando cola %s: %v", queueName, err)
	}
	log.Printf("Conexión a RabbitMQ establecida, cola: %s", queueName)

	// Iniciar servidor gRPC
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpServer := grpc.NewServer()
	pb.RegisterWeatherServiceServer(grpServer, &server{
		rabbitConn: rabbitConn,
		rabbitCh:   rabbitCh,
		queueName:  queueName,
	})
	log.Printf("Servidor gRPC RabbitMQ corriendo en :50052")
	if err := grpServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}