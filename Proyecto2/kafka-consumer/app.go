package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

// WeatherData refleja la estructura de los mensajes enviados
type WeatherData struct {
	Description string `json:"description"`
	Country     string `json:"country"`
	Weather     string `json:"weather"`
}

func main() {
	// Configuración desde variables de entorno
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "my-topic"
	}
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "my-cluster-kafka-bootstrap.kafka:9092"
	}
	groupID := os.Getenv("KAFKA_GROUP_ID")
	if groupID == "" {
		groupID = "go-consumer-group"
	}

	// Configurar lector de Kafka
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: groupID,
	})
	defer r.Close()

	log.Printf("Consumidor Kafka iniciado para topic %s, broker %s", topic, broker)

	// Inicar sesion en Redis

	ctx := context.Background()

    redisHost := os.Getenv("REDIS_HOST")
    redisPassword := os.Getenv("REDIS_PASSWORD")

    rdb := redis.NewClient(&redis.Options{
        Addr:     redisHost,
        Password: redisPassword,
        DB:       0,
    })


	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error leyendo mensaje: %v", err)
			continue
		}

		// Deserializar el mensaje
		var data WeatherData
		if err := json.Unmarshal(m.Value, &data); err != nil {
			log.Printf("Error deserializando mensaje: %v", err)
			continue
		}

        timestamp := time.Now().Format("20060102_150405")
        key := fmt.Sprintf("registro:%s", timestamp)

		// Grabar mensaje en redis

        err = rdb.Set(ctx, key, m.Value, 0).Err()
        if err != nil {
            log.Printf("Error al guardar en Redis: %v\n", err)
        } else {
            log.Printf("Registro guardado [%s]: %s\n", key, m.Value)
        }



		log.Printf("Mensaje recibido: offset=%d, partition=%d, value=%+v, y guardado en Redis", m.Offset, m.Partition, data)
	}
}