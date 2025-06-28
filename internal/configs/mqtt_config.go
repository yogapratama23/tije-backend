package configs

import (
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func CreateMqttClient() mqtt.Client {
	mqttOppt := mqtt.NewClientOptions().AddBroker(os.Getenv("MQTT_BROKER_URL"))
	mqttOppt.SetClientID("tije-mqtt")
	mqttOppt.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
	})
	mqttClient := mqtt.NewClient(mqttOppt)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return mqttClient
}
