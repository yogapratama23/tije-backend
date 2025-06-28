package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yogapratama23/tije-backend/internal/dtos"
	"github.com/yogapratama23/tije-backend/internal/services"
)

type VehicleLocationController interface {
	Insert(ctx *gin.Context)
	LatestLocation(ctx *gin.Context)
	History(ctx *gin.Context)
	SimulateDriving(ctx *gin.Context)
}

type vehicleLocationController struct {
	vehicleLocationService services.VehicleLocationService
}

func NewVehicleLocationController(
	vehicleLocationService services.VehicleLocationService,
) VehicleLocationController {
	return &vehicleLocationController{vehicleLocationService: vehicleLocationService}
}

func (c *vehicleLocationController) Insert(ctx *gin.Context) {
	var err error
	var json dtos.CreateVehicleLocationInput
	if err = ctx.ShouldBindJSON(&json); err != nil {
		log.Printf("[VehicleLocationController.Insert] error: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid parameters",
		})
		return
	}

	validate := json.Validate()
	if validate != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": validate,
		})
		return
	}

	err = c.vehicleLocationService.Create(ctx.Request.Context(), json)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Something went wrong",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Vehicle Location Created",
	})
}

func (c *vehicleLocationController) LatestLocation(ctx *gin.Context) {
	var err error
	var json dtos.FindOneVehicleLocationFilter
	if err = ctx.ShouldBindUri(&json); err != nil {
		log.Printf("[VehicleLocationController.LatestLocation] error: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid parameters",
		})
		return
	}

	latest := true
	json.Latest = &latest

	vl, err := c.vehicleLocationService.FindOne(ctx.Request.Context(), json)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Something went wrong",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": vl,
	})
}

func (c *vehicleLocationController) History(ctx *gin.Context) {
	var err error
	var json dtos.FindManyVehicleLocationFilter
	if err = ctx.ShouldBindQuery(&json); err != nil {
		log.Printf("[VehicleLocationController.History] error: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid parameters",
		})
		return
	}

	validate := json.Validate()
	if validate != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": validate,
		})
		return
	}

	latest := true
	json.Latest = &latest

	vehicleLocations, err := c.vehicleLocationService.FindMany(ctx.Request.Context(), json)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Something went wrong",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": vehicleLocations,
	})
}

func (c *vehicleLocationController) SimulateDriving(ctx *gin.Context) {
	vehicleId := ctx.Query("vehicle_id")
	if vehicleId == "" {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Vehicle id is required",
		})
		return
	}
	c.vehicleLocationService.SimulateDriving(vehicleId)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "You have arrive",
	})
}
