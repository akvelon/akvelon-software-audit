package main

import (
	"akvelon/akvelon-software-audit/license-audit-service/pkg/licanalize"
	"akvelon/akvelon-software-audit/license-audit-service/pkg/storage/mongo"
	"log"

	"github.com/streadway/amqp"
)

const (
	uXAuditQueueName = "audit-queue"
	rabbitSrv        = "amqp://guest:guest@rabbitmq:5672"
)

func main() {
	s := new(mongo.Storage)
	err := s.InitStorage()
	if err != nil {
		log.Fatal("ERROR: could not init storage: ", err)
	}

	licAnalizer := licanalize.NewService(s)

	conn, err := amqp.Dial(rabbitSrv)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	defer ch.Close()

	q, err := ch.QueueDeclare(
		uXAuditQueueName, // name
		false,            // durable
		false,            // delete when usused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)

	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	forever := make(chan bool)
	log.Printf("Audit worker is running and listening to '%s'", q.Name)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			err := licAnalizer.Scan(licanalize.AnalizedRepo{
				URL: string(d.Body),
			})
			if err != nil {
				log.Println(err)
				log.Printf("Failed to Scan repo %s", d.Body)
			}
		}
	}()

	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
