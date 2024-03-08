package nsq

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"log"
	CO "sequency/config"
	M "sequency/models"
)

func SendMessageToNSQ(topic string, message M.MessageNSQ) {
	config := nsq.NewConfig()
	envs := CO.ConfigEnv()
	producer, err := nsq.NewProducer(envs["NSQ"], config)

	if err != nil {
		log.Fatal(err)
	}

	payload, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	err = producer.Publish(topic, payload)

	if err != nil {
		log.Println(err)
	}

	defer producer.Stop()

}

func ProcessOrderNSQ() {
	config := nsq.NewConfig()
	envs := CO.ConfigEnv()
	consumer, err := nsq.NewConsumer("emails", "email_chanels", config)

	if err != nil {
		log.Fatal(err)
	}

	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Println(message)
		return nil
	}))

	err = consumer.ConnectToNSQD(envs["NSQ"])

	if err != nil {
		log.Println(err)
	}

	<-make(chan struct{})

}
