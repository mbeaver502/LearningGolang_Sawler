package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // name of the exchange
		"topic",      // kind/type of exchange
		true,         // exchange is durable
		false,        // do not auto-delete
		false,        // this exchange is not internal only
		false,        // no-wait
		nil,          // no specific arguments
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name
		false, // NOT durable
		false, // do NOT auto-delete
		true,  // this is exclusive
		false, // no-wait
		nil,   // no specific arguments
	)
}
