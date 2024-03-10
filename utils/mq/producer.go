package mq

import (
	"context"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type ConnectionMQ struct {
	MQ *amqp091.Connection
}

func MQ() *amqp091.Connection {
	fmt.Println("RabbitMQ in Golang")

	connection, err := amqp091.Dial("amqp://sqcgo:sqcgo@localhost:5672/sqcgo")
	if err != nil {
		panic(err)
	}
	// defer connection.Close()

	fmt.Println("Successfully connected to RabbitMQ instance")

	ch, err := connection.Channel()

	if err != nil {
		panic(err)
	}

	defer ch.Close()

	_, err = ch.QueueDeclare(
		"testing", // nombre de la cola
		true,      // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		panic(err)
	}

	return connection
}

func (C ConnectionMQ) SendMessage(body []byte) {

	channel, err := C.MQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"testing", // name
		true,      // durable
		false,     // auto delete
		false,     // exclusive
		false,     // no wait
		nil,       // args
	)

	if err != nil {
		panic(err)
	}

	err = channel.PublishWithContext(
		context.TODO(),
		"",        // exchange
		"testing", // key
		false,     // mandatory
		false,     // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("Queue status:", queue)
	fmt.Println("Successfully published message")
}

func (C ConnectionMQ) PollMq() {

	channel, err := C.MQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()
	msgs, err := channel.Consume(
		"testing",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			fmt.Printf("Received Message: %s\n", msg.Body)
		}
	}()

	fmt.Println("Waiting for messages...")
	<-forever
}
