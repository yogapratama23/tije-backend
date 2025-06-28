package dtos

type CreateVehicleLocationInput struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

func (d *CreateVehicleLocationInput) Validate() string {
	if d.VehicleID == "" {
		return "Vehicle ID is required"
	}
	if d.Latitude == 0 {
		return "Latitude is required"
	}
	if d.Longitude == 0 {
		return "Longitude is required"
	}
	if d.Timestamp == 0 {
		return "Timestamp is required"
	}

	return ""
}

type FindOneVehicleLocationFilter struct {
	ID        *string `uri:"id" json:"id"`
	VehicleID *string `uri:"vehicle_id" json:"vehicle_id"`
	Latest    *bool
}

type FindManyVehicleLocationFilter struct {
	Start  *int64 `form:"start" json:"start"`
	End    *int64 `form:"end" json:"end"`
	Latest *bool
}

func (d *FindManyVehicleLocationFilter) Validate() string {
	if d.Start == nil && d.End == nil {
		return ""
	}
	if d.Start == nil && d.End != nil {
		return "Start time is required"
	}
	if d.Start != nil && d.End == nil {
		return "End time is required"
	}
	if d.Start != nil && d.End != nil {
		if *d.Start > *d.End {
			return "Start time can't be greater than end time"
		}
	}

	return ""
}

type GeofenceAlertLocationType struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type GeofenceAlertInput struct {
	VehicleID string                    `json:"vehicle_id"`
	Event     string                    `json:"event"`
	Location  GeofenceAlertLocationType `json:"location"`
	Timestamp int64                     `json:"timestamp"`
}
