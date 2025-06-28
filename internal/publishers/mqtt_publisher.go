package publishers

import mqtt "github.com/eclipse/paho.mqtt.golang"

type MqttPublisher interface {
	Publish(topic string, msg any) error
}

type mqttPublisher struct {
	client mqtt.Client
}

func NewMqttPublisher(
	client mqtt.Client,
) MqttPublisher {
	return &mqttPublisher{client: client}
}

func (p *mqttPublisher) Publish(topic string, msg any) error {
	token := p.client.Publish(topic, 0, false, msg)
	token.Wait()

	return token.Error()
}
