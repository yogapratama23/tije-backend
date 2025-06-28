package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yogapratama23/tije-backend/internal/configs"
	"github.com/yogapratama23/tije-backend/internal/controllers"
	"github.com/yogapratama23/tije-backend/internal/database"
	"github.com/yogapratama23/tije-backend/internal/database/repositories"
	"github.com/yogapratama23/tije-backend/internal/publishers"
	"github.com/yogapratama23/tije-backend/internal/services"
	"github.com/yogapratama23/tije-backend/internal/subscribers"
)

func main() {
	router := gin.Default()
	mqttClient := configs.CreateMqttClient()

	rabbitMqConn := configs.CreateRabbitMqConn()
	defer rabbitMqConn.Close()

	rabbitMqChan, err := rabbitMqConn.Channel()
	if err != nil {
		panic(err.Error())
	}
	defer rabbitMqChan.Close()

	configs.SetupRabbitMq(rabbitMqChan)

	db := database.Connect()
	defer db.Close(context.Background())

	database.RunMigration(db)

	// repositories
	vehicleLocationRepository := repositories.NewVehicleLocationRepository(db)

	// publishers
	mqttPublisher := publishers.NewMqttPublisher(mqttClient)
	rabbitMqProducer := publishers.NewRabbitMqProducer(rabbitMqChan)

	// services
	vehicleLocationService := services.NewVehicleLocationService(vehicleLocationRepository, rabbitMqProducer, mqttPublisher)

	// controllers
	vehicleLocationController := controllers.NewVehicleLocationController(vehicleLocationService)

	// subscribers
	mqttSubscriber := subscribers.NewMqttSubscriber(vehicleLocationService)

	// routes
	vehicleLocationRoute := router.Group("/vehicles")

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})

	vehicleLocationRoute.POST("/create", vehicleLocationController.Insert)
	vehicleLocationRoute.GET("/:vehicle_id/location", vehicleLocationController.LatestLocation)
	vehicleLocationRoute.GET("/:vehicle_id/history", vehicleLocationController.History)

	router.GET("/simulate-driving", vehicleLocationController.SimulateDriving)

	mqttClient.Subscribe("/fleet/vehicle/+/location", 0, mqttSubscriber.VehicleLocationSubscriber)

	router.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))) // listen and serve on 0.0.0.0:8080
}
