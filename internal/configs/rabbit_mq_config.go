package configs

import (
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeName = "fleet.events"
	QueueName    = "geofence_alert"
	RoutingKey   = "my_key"
)

func CreateRabbitMqConn() *amqp.Connection {
	var conn *amqp.Connection
	var err error

	for i := 1; i < 21; i++ {
		conn, err = amqp.Dial(os.Getenv("RABBIT_MQ_URL"))
		if err == nil {
			fmt.Printf("Connected to RabbitMQ on attempt number %d", i)
			return conn
		}

		fmt.Printf("Connection failed on attempt %d \n", i)
		time.Sleep(2 * time.Second)
	}

	panic("Failed to connect to RabbitMQ \n")
}

func SetupRabbitMq(rabbitMqChan *amqp.Channel) {
	if err := rabbitMqChan.ExchangeDeclare(ExchangeName, "direct", true, false, false, false, nil); err != nil {
		panic(err.Error())
	}
	queue, err := rabbitMqChan.QueueDeclare(QueueName, false, false, false, false, nil)
	if err != nil {
		panic(err.Error())
	}
	if err := rabbitMqChan.QueueBind(queue.Name, RoutingKey, ExchangeName, false, nil); err != nil {
		panic(err.Error())
	}
}
