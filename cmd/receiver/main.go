package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/yogapratama23/tije-backend/internal/configs"
)

func main() {
	rabbitMqConn := configs.CreateRabbitMqConn()
	defer rabbitMqConn.Close()

	rabbitMqChan, err := rabbitMqConn.Channel()
	if err != nil {
		panic(err.Error())
	}
	defer rabbitMqChan.Close()

	q, err := rabbitMqChan.QueueDeclare(configs.QueueName, false, false, false, false, nil)
	if err != nil {
		panic(err.Error())
	}

	msgs, err := rabbitMqChan.Consume(q.Name, "", true, false, false, false, nil)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
