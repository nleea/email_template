package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	EM "sequency/utils/emails"
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

type EmailInfoData struct {
	Subject      string `json:"subject"`
	To_address   string `json:"to_address"`
	From_address string `json:"fom_address"`
	From_name    string `json:"from_name"`
}

type SendEmailStrunc struct {
	Email_info []map[string]string `json:"email_info"`
	Data       any                 `json:"data"`
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
			var dataAggregation []SendEmailStrunc
			json.Unmarshal(msg.Body, &dataAggregation)

			EM.SendEmail(
				dataAggregation[0].Email_info[2]["Value"],
				dataAggregation[0].Email_info[1]["Value"],
				dataAggregation[0].Email_info[0]["Value"],
				"Text")

		}
	}()

	fmt.Println("Waiting for messages...")
	<-forever
}
