package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

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

		// Contador total de mensajes
		err = rdb.Incr(ctx, "total:messages").Err()
		if err != nil {
			log.Printf("Error incrementando contador total: %v", err)
		}

		// Contador por país (usando hash)
		err = rdb.HIncrBy(ctx, "countries:count", data.Country, 1).Err()
		if err != nil {
			log.Printf("Error incrementando contador para país %s: %v", data.Country, err)
		}



		log.Printf("Mensaje recibido: offset=%d, partition=%d, value=%+v, y guardado en Redis", m.Offset, m.Partition, data)
	}
}