package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer models the events received from the queue.
type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

// NewConsumer will set up a new RabbitMQ/AMQP consumer.
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		log.Println(err)
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		log.Println(err)
		return err
	}

	return declareExchange(channel)
}

// Payload models what will be pushed to the queue.
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// Listen listens for given topics.
func (consumer *Consumer) Listen(topics []string) error {
	// get a channel on our connection
	ch, err := consumer.conn.Channel()
	if err != nil {
		log.Println(err)
		return err
	}
	defer ch.Close()

	// create a queue on our channel
	q, err := declareRandomQueue(ch)
	if err != nil {
		log.Println(err)
		return err
	}

	// bind topics to the queue on our channel
	for _, s := range topics {
		err := ch.QueueBind(
			q.Name,       // name
			s,            // key
			"logs_topic", // exchange
			false,        // no-wait
			nil,          // no specific arguments
		)

		if err != nil {
			log.Println(err)
			return err
		}
	}

	// set up a chan to start consume messages
	messages, err := ch.Consume(
		q.Name, // name
		"",     // consumer
		true,   // auto-acknowledge
		false,  // exclusive?
		false,  // no-local
		false,  // no-wait
		nil,    // no specific arguments
	)

	if err != nil {
		log.Println(err)
		return err
	}

	// set up a chan so we can listen forever
	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	log.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		// authenticate
		break
	// we can have as many cases as we want,
	// as long as we write the logic for them
	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(entry Payload) error {
	// turn the log request into JSON we can send to Logger
	jsonData, _ := json.Marshal(entry)

	// call the service
	logServiceURL := "http://logger-service/log" // logger-service is the name inside Docker
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return err
	}
	defer response.Body.Close()

	// ensure we got the right response status code
	if response.StatusCode != http.StatusAccepted {
		log.Printf("invalid status code: %d", response.StatusCode)
		return fmt.Errorf("invalid status code: %d", response.StatusCode)
	}

	return nil
}
