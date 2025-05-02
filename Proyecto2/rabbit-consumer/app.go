package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type WeatherData struct {
	Description string `json:"description"`
	Country     string `json:"country"`
	Weather     string `json:"weather"`
}

func main() {
	// Configuración de RabbitMQ desde variables de entorno
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@rabbitmq.rabbitmq.svc.cluster.local:5672/"
	}
	queueName := os.Getenv("RABBITMQ_QUEUE")
	if queueName == "" {
		queueName = "weather-queue"
	}

	// Conectar a RabbitMQ
	var conn *amqp091.Connection
	var err error
	for i := 0; i < 5; i++ {
		conn, err = amqp091.Dial(rabbitMQURL)
		if err == nil {
			break
		}
		log.Printf("Error conectando a RabbitMQ (intento %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if conn == nil {
		log.Fatalf("No se pudo conectar a RabbitMQ después de varios intentos")
	}
	defer conn.Close()

	// Crear canal
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creando canal de RabbitMQ: %v", err)
	}
	defer ch.Close()

	// Declarar cola
	q, err := ch.QueueDeclare(
		queueName, // nombre de la cola
		true,      // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // noWait
		nil,       // argumentos
	)
	if err != nil {
		log.Fatalf("Error declarando cola %s: %v", queueName, err)
	}

	// Consumir mensajes
	msgs, err := ch.Consume(
		q.Name, // cola
		"",     // consumidor
		true,   // autoAck
		false,  // exclusive
		false,  // noLocal
		false,  // noWait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Error registrando consumidor: %v", err)
	}

	log.Printf("Consumidor RabbitMQ iniciado para cola %s", queueName)
	for msg := range msgs {
		var weatherData WeatherData
		err := json.Unmarshal(msg.Body, &weatherData)
		if err != nil {
			log.Printf("Error deserializando mensaje: %v", err)
			continue
		}
		log.Printf("Mensaje recibido: Description=%s, Country=%s, Weather=%s",
			weatherData.Description, weatherData.Country, weatherData.Weather)
	}
}