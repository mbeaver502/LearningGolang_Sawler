package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	log.Printf("Broker listening on port %s\n", webPort)

	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}

	// define an HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
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
		log.Printf("backing off for %v\n", backoff)
		time.Sleep(backoff)
	}

	log.Println("*** Connected to RabbitMQ! ***")

	return connection, nil
}
