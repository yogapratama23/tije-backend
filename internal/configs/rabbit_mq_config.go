package configs

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeName = "fleet.events"
	QueueName    = "geofence_alert"
	RoutingKey   = "my_key"
)

func CreateRabbitMqConn() *amqp.Connection {
	rabbitMqConn, err := amqp.Dial(os.Getenv("RABBIT_MQ_URL"))
	if err != nil {
		panic(err.Error())
	}

	return rabbitMqConn
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
