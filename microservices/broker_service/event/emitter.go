package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp.Connection
}

// NewEventEmitter sets up a new event emitter for RabbitMQ.
func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		log.Println(err)
		return Emitter{}, err
	}

	return emitter, nil
}

func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		log.Println(err)
		return err
	}
	defer channel.Close()

	return declareExchange(channel)
}

// Push pushes an event onto a RabbitMQ connection channel.
func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		log.Println(err)
		return err
	}
	defer channel.Close()

	log.Println("Pushing to channel")

	// Note: channel.Publish is deprecated
	err = channel.Publish(
		"logs_topic", // exchange
		severity,     // key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}