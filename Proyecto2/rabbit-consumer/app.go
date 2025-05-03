package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type WeatherData struct {
	Description string `json:"description"`
	Country     string `json:"country"`
	Weather     string `json:"weather"`
}

func main() {
	// Configuración desde variables de entorno
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@rabbitmq.rabbitmq.svc.cluster.local:5672/"
	}
	
	queueName := os.Getenv("RABBITMQ_QUEUE")
	if queueName == "" {
		queueName = "weather-queue"
	}
	
	valkeyHost := os.Getenv("VALKEY_HOST")
	if valkeyHost == "" {
		valkeyHost = "valkey-service.proyecto2.svc.cluster.local:6379"
	}
	
	valkeyPassword := os.Getenv("VALKEY_PASSWORD")

	// Inicializar cliente Valkey (compatible con Redis)
	valkeyClient := redis.NewClient(&redis.Options{
		Addr:     valkeyHost,
		Password: valkeyPassword,
		DB:       0, // Base de datos por defecto
	})

	// Verificar conexión a Valkey
	ctx := context.Background()
	_, err := valkeyClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error conectando a Valkey: %v", err)
	}
	log.Println("Conexión a Valkey establecida")

	// Configurar conexión a RabbitMQ con reintentos
	var rabbitConn *amqp091.Connection
	maxRetries := 5
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		rabbitConn, err = amqp091.Dial(rabbitMQURL)
		if err == nil {
			break
		}
		
		log.Printf("Error conectando a RabbitMQ (intento %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}

	if rabbitConn == nil {
		log.Fatalf("No se pudo conectar a RabbitMQ después de %d intentos", maxRetries)
	}
	defer rabbitConn.Close()
	log.Println("Conexión a RabbitMQ establecida")

	// Crear canal RabbitMQ
	channel, err := rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Error creando canal RabbitMQ: %v", err)
	}
	defer channel.Close()

	// Declarar cola durable
	queue, err := channel.QueueDeclare(
		queueName, // nombre
		true,      // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // noWait
		nil,       // argumentos
	)
	if err != nil {
		log.Fatalf("Error declarando cola: %v", err)
	}

	// Configurar consumidor
	messages, err := channel.Consume(
		queue.Name, // cola
		"",         // consumer tag
		false,     // auto-ack (manualmente para evitar pérdidas)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("Error registrando consumidor: %v", err)
	}

	log.Printf("Consumidor iniciado para cola: %s", queueName)

	// Procesar mensajes
	for msg := range messages {
		startTime := time.Now()
		
		var weather WeatherData
		if err := json.Unmarshal(msg.Body, &weather); err != nil {
			log.Printf("Error deserializando mensaje (ID: %s): %v", msg.MessageId, err)
			msg.Nack(false, true) // Reintentar mensaje
			continue
		}

		// Contador total de mensajes
		err = valkeyClient.Incr(ctx, "total:messages").Err()
		if err != nil {
			log.Printf("Error incrementando contador total: %v", err)
		}

		// Contador por país (usando hash)
		err = valkeyClient.HIncrBy(ctx, "countries:count", weather.Country, 1).Err()
		if err != nil {
			log.Printf("Error incrementando contador para país %s: %v", weather.Country, err)
		}

		// Registrar éxito
		processingTime := time.Since(startTime).Milliseconds()

		log.Printf("Mensaje recibido: Description=%s, Country=%s, Weather=%s, tiempo de procesamiento=%d ms",
		weather.Description, weather.Country, weather.Weather, processingTime)

		msg.Ack(false)
	}
}