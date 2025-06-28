package subscribers

import (
	"context"
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/yogapratama23/tije-backend/internal/dtos"
	"github.com/yogapratama23/tije-backend/internal/services"
)

type MqttSubscriber interface {
	VehicleLocationSubscriber(client mqtt.Client, msg mqtt.Message)
}

type mqttSubscriber struct {
	vehicleLocationService services.VehicleLocationService
}

func NewMqttSubscriber(
	vehicleLocationService services.VehicleLocationService,
) MqttSubscriber {
	return &mqttSubscriber{vehicleLocationService: vehicleLocationService}
}

func (s *mqttSubscriber) VehicleLocationSubscriber(client mqtt.Client, msg mqtt.Message) {
	var input dtos.CreateVehicleLocationInput
	err := json.Unmarshal(msg.Payload(), &input)
	if err != nil {
		log.Printf("[MqttSubscriber.VehicleLocationSubscriber.Unmarshall] error: %s", err.Error())
	}

	validate := input.Validate()
	if validate != "" {
		log.Printf("[MqttSubscriber.VehicleLocationSubscriber.Validation] error: %s", validate)
	}

	err = s.vehicleLocationService.Create(context.Background(), input)
	if err != nil {
		log.Printf("[MqttSubscriber.VehicleLocationSubscriber.Create] failed to create vehicle location for %s", input.VehicleID)
	}
}
