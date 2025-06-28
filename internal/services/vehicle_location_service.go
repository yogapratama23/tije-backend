package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/yogapratama23/tije-backend/internal/configs"
	"github.com/yogapratama23/tije-backend/internal/database/models"
	"github.com/yogapratama23/tije-backend/internal/database/repositories"
	"github.com/yogapratama23/tije-backend/internal/dtos"
	"github.com/yogapratama23/tije-backend/internal/publishers"
)

type VehicleLocationService interface {
	Create(ctx context.Context, input dtos.CreateVehicleLocationInput) error
	FindOne(ctx context.Context, input dtos.FindOneVehicleLocationFilter) (*models.VehicleLocation, error)
	FindMany(ctx context.Context, input dtos.FindManyVehicleLocationFilter) ([]*models.VehicleLocation, error)
	SimulateDriving(vehicleId string)
}

type vehicleLocationService struct {
	vehicleLocationRepository repositories.VehicleLocationRepository
	rabbitMqProducer          publishers.RabbitMqProducer
	mqttPublisher             publishers.MqttPublisher
}

func NewVehicleLocationService(
	vehicleLocationRepository repositories.VehicleLocationRepository,
	rabbitMqProduceer publishers.RabbitMqProducer,
	mqttPublisher publishers.MqttPublisher,
) VehicleLocationService {
	return &vehicleLocationService{
		vehicleLocationRepository: vehicleLocationRepository,
		rabbitMqProducer:          rabbitMqProduceer,
		mqttPublisher:             mqttPublisher,
	}
}

func (s *vehicleLocationService) Create(ctx context.Context, input dtos.CreateVehicleLocationInput) error {
	_, err := s.vehicleLocationRepository.Insert(ctx, input)
	if err != nil {
		log.Printf("[VehicleLocationService.Create] error: %s", err.Error())
		return err
	}

	return nil
}

func (s *vehicleLocationService) FindOne(ctx context.Context, input dtos.FindOneVehicleLocationFilter) (*models.VehicleLocation, error) {
	vl, err := s.vehicleLocationRepository.FindOne(ctx, input)
	if err != nil {
		log.Printf("[VehicleLocationService.FindOne] error %s", err.Error())
	}

	return vl, err
}

func (s *vehicleLocationService) SimulateDriving(vehicleId string) {
	fmt.Println("Start Driving")

	start := dtos.GeofenceAlertLocationType{
		Latitude:  -6.1936,
		Longitude: 106.8200,
	}
	end := dtos.GeofenceAlertLocationType{
		Latitude:  -6.1700,
		Longitude: 106.8249,
	}
	checkPoints := []dtos.GeofenceAlertLocationType{
		{
			Latitude:  -6.2000,
			Longitude: 106.8160,
		},
		{
			Latitude:  -6.2146,
			Longitude: 106.8451,
		},
		{
			Latitude:  -6.1702,
			Longitude: 106.8249,
		},
	}
	// hitted checkpoint
	checkPoint := dtos.GeofenceAlertLocationType{
		Latitude:  -6.2000,
		Longitude: 106.816,
	}

	steps := 20
	radius := 50.0
	path := []dtos.GeofenceAlertLocationType{}
	stepsToCheckpoint := steps / 2
	stepsToEnd := steps - stepsToCheckpoint

	for i := 0; i < stepsToCheckpoint; i++ {
		lat := start.Latitude + (checkPoint.Latitude-start.Latitude)*float64(i)/float64(stepsToCheckpoint)
		lng := start.Longitude + (checkPoint.Longitude-start.Longitude)*float64(i)/float64(stepsToCheckpoint)
		path = append(path, dtos.GeofenceAlertLocationType{
			Latitude:  lat,
			Longitude: lng,
		})
	}

	path = append(path, checkPoint)

	for i := 1; i <= stepsToEnd; i++ {
		lat := checkPoint.Latitude + (end.Latitude-checkPoint.Latitude)*float64(i)/float64(stepsToEnd)
		lng := checkPoint.Longitude + (end.Longitude-checkPoint.Longitude)*float64(i)/float64(stepsToEnd)
		path = append(path, dtos.GeofenceAlertLocationType{
			Latitude:  lat,
			Longitude: lng,
		})
	}

	for i, pos := range path {
		fmt.Printf("Step %d - Location: (%.6f, %.6f) \n", i, pos.Latitude, pos.Longitude)

		mqttInput := dtos.CreateVehicleLocationInput{
			VehicleID: vehicleId,
			Latitude:  pos.Latitude,
			Longitude: pos.Longitude,
			Timestamp: time.Now().Unix(),
		}

		marshalInput, err := json.Marshal(mqttInput)
		if err != nil {
			fmt.Printf("Failed to marshal json for %v", pos)
			continue
		}

		err = s.mqttPublisher.Publish(fmt.Sprintf("/fleet/vehicle/%s/location", vehicleId), marshalInput)
		if err != nil {
			fmt.Printf("Failed to publish message to mqtt broker: %s", err.Error())
		}

		for _, cp := range checkPoints {
			distance := calculateDistance(pos.Latitude, pos.Longitude, cp.Latitude, cp.Longitude)
			if distance <= radius {
				fmt.Printf("Within %.0f meters of checkpoint at (%.4f, %.4f) \n", distance, cp.Latitude, cp.Longitude)
				geofenceInput := dtos.GeofenceAlertInput{
					VehicleID: vehicleId,
					Event:     "geofence_entry",
					Location: dtos.GeofenceAlertLocationType{
						Latitude:  pos.Latitude,
						Longitude: pos.Longitude,
					},
					Timestamp: time.Now().Unix(),
				}

				s.rabbitMqProducer.Publish(configs.QueueName, geofenceInput)
			}
		}

		time.Sleep(2 * time.Second)
	}

	fmt.Println("Finish Driving")
}

func (s *vehicleLocationService) FindMany(ctx context.Context, input dtos.FindManyVehicleLocationFilter) ([]*models.VehicleLocation, error) {
	vehicleLocations, err := s.vehicleLocationRepository.FindMany(ctx, input)
	if err != nil {
		log.Printf("[VehicleLocationService.FindMany] error: %s", err.Error())
	}

	return vehicleLocations, err
}

func calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371e3
	phi1 := lat1 * math.Pi / 180
	phi2 := lat2 * math.Pi / 180
	deltaPhi := (lat2 - lat1) * math.Pi / 180
	deltaLambda := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) + math.Cos(phi1)*math.Cos(phi2)*math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
