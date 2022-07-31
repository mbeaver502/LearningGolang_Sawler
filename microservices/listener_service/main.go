package main

import (
	"errors"
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

	log.Println("*** Connected to RabbitMQ! ***")

	// start listening for messages

	// create consumer

	// watch the queue and consume events
}

func connect() (*amqp.Connection, error) {
	var counts int64
	backoff := 1 * time.Second

	var connection *amqp.Connection

	// don't continue until rabbitmq is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost")
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

	return connection, nil
}
