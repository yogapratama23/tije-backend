package publishers

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	GeofenceEntry = "geofence_entry"
)

type RabbitMqProducer interface {
	Publish(queueName string, body any)
}

type rabbitMqProducer struct {
	rabbitMqChan *amqp.Channel
}

func NewRabbitMqProducer(
	rabbitMqChan *amqp.Channel,
) RabbitMqProducer {
	return &rabbitMqProducer{rabbitMqChan: rabbitMqChan}
}

func (p *rabbitMqProducer) Publish(queueName string, body any) {
	input, err := json.Marshal(body)
	if err != nil {
		panic(err.Error())
	}

	err = p.rabbitMqChan.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        input,
	})
	if err != nil {
		panic(err.Error())
	}
}
