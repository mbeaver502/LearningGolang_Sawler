package main

import (
	"errors"
	"listenerservice/event"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	log.Println("Listener started")

	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer rabbitConn.Close()

	// start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		log.Panicln(err)
	}

	// watch the queue and consume events
	topics := []string{"log.INFO", "log.WARNING", "log.ERROR"}
	err = consumer.Listen(topics)
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	backoff := 1 * time.Second

	var connection *amqp.Connection

	// don't continue until rabbitmq is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq") // rabbitmq is the name of the service inside Docker
		if err != nil {
			log.Println("rabbitmq not yet ready")
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			log.Println("failed to connect after 5 tries")
			return nil, errors.New("failed to connect after 5 tries")
		}

		// exponential backoff
		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Printf("backing off for %v seconds\n", backoff)
		time.Sleep(backoff)
	}

	log.Println("*** Connected to RabbitMQ! ***")

	return connection, nil
}
