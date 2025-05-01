package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	pb "grpc-cliente/proto"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WeatherData struct {
	Description string `json:"description"`
	Country     string `json:"country"`
	Weather     string `json:"weather"`
}

func handleWeather(w http.ResponseWriter, r *http.Request) {
	var data WeatherData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error al leer JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Recibido desde Rust: %+v\n", data)

	// Llamar al cliente gRPC
	go func() {
		err := callGrpcClient(data)
		if err != nil {
			log.Printf("Error llamando a gRPC: %v", err)
		}
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Datos recibidos correctamente"))
}


func callGrpcClient(data WeatherData) error {
	// Conectar al servidor grpc-server-rabit
	conn, err := grpc.Dial("grpc-server-rabbit:50052", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(), grpc.WithTimeout(3*time.Second))
	if err != nil {
		return err
	}
	defer conn.Close()

	// Crear el clima a partir del stub generado
	client := pb.NewWeatherServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Convertir WeatherData a gRPC WeatherData
	grpcData := &pb.WeatherData{
		Description: data.Description,
		Country:     data.Country,
		Weather:     data.Weather,
	}

	// Llamar a PostToRabbitMQ
	rabbitResp, err := client.PostToRabbitMQ(ctx, grpcData)
	if err != nil {
		log.Printf("Error en PostToRabbitMQ: %v", err)
		return err
	}
	log.Printf("Respuesta RabbitMQ: %s", rabbitResp.Message)


	// Conectar al servidor grpc-server-rabit
	conn, err = grpc.Dial("grpc-server-kafka:50051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(), grpc.WithTimeout(3*time.Second))
	if err != nil {
		return err
	}
	defer conn.Close()

	// Crear el clima a partir del stub generado
	client = pb.NewWeatherServiceClient(conn)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()


	// Llamar a PostToKafka
	kafkaResp, err := client.PostToKafka(ctx, grpcData)
	if err != nil {
		log.Printf("Error en PostToKafka: %v", err)
		return err
	}
	log.Printf("Respuesta Kafka: %s", kafkaResp.Message)

	return nil
}



func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handleWeather).Methods("POST")
	log.Println("Servidor REST corriendo en :8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatalf("Error al iniciar el servidor REST: %v", err)
	}
}
