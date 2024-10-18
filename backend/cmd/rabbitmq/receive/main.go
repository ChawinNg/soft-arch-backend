package main

import (
	"log"

	"github.com/streadway/amqp"
)

func connectRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://root:root@localhost:5672/")
	if err != nil {
		return nil, nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}
	return conn, ch, nil
}

func receiveMessages(ch *amqp.Channel, queueName string) {
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %s", err)
	}

	log.Printf("[*] Waiting for Messages.")

	for msg := range msgs {
		log.Printf("Received message: %s", msg.Body)
	}
}

func main() {
	conn, ch, err := connectRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()
	defer ch.Close()

	receiveMessages(ch, "enrollment")
}
