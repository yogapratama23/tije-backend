package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/yogapratama23/tije-backend/internal/database/models"
	"github.com/yogapratama23/tije-backend/internal/dtos"
)

type VehicleLocationRepository interface {
	Insert(ctx context.Context, input dtos.CreateVehicleLocationInput) (*int, error)
	FindOne(ctx context.Context, filter dtos.FindOneVehicleLocationFilter) (*models.VehicleLocation, error)
	FindMany(ctx context.Context, filter dtos.FindManyVehicleLocationFilter) ([]*models.VehicleLocation, error)
}

type vehicleLocationRepository struct {
	db *pgx.Conn
}

func NewVehicleLocationRepository(db *pgx.Conn) VehicleLocationRepository {
	return &vehicleLocationRepository{db: db}
}

func (r *vehicleLocationRepository) Insert(ctx context.Context, input dtos.CreateVehicleLocationInput) (*int, error) {
	var err error
	query := `
		INSERT INTO vehicle_locations (vehicle_id, latitude, longitude, timestamp)
		VALUES (@vehicle_id, @latitude, @longitude, @timestamp)
		RETURNING id
	`
	queryArgs := pgx.NamedArgs{
		"vehicle_id": input.VehicleID,
		"latitude":   input.Latitude,
		"longitude":  input.Longitude,
		"timestamp":  input.Timestamp,
	}

	var id int
	err = r.db.QueryRow(ctx, query, queryArgs).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &id, nil
}

func (r *vehicleLocationRepository) FindOne(ctx context.Context, filter dtos.FindOneVehicleLocationFilter) (*models.VehicleLocation, error) {
	var err error
	vehicleLocation := models.VehicleLocation{}

	query := `
		SELECT id, vehicle_id, latitude, longitude, timestamp
		FROM vehicle_locations
		WHERE 1=1
	`
	queryArgs := pgx.NamedArgs{}

	if filter.ID != nil {
		query += " AND id = @id"
		queryArgs["id"] = filter.ID
	}
	if filter.VehicleID != nil {
		query += " AND vehicle_id = @vehicle_id"
		queryArgs["vehicle_id"] = filter.VehicleID
	}

	if filter.Latest != nil && *filter.Latest {
		query += " ORDER BY timestamp DESC"
	}

	err = r.db.QueryRow(ctx, query, queryArgs).Scan(&vehicleLocation.ID, &vehicleLocation.VehicleID, &vehicleLocation.Latitude, &vehicleLocation.Longitude, &vehicleLocation.Timestamp)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &vehicleLocation, nil
}

func (r *vehicleLocationRepository) FindMany(ctx context.Context, filter dtos.FindManyVehicleLocationFilter) ([]*models.VehicleLocation, error) {
	vehicleLocations := make([]*models.VehicleLocation, 0)
	var err error

	query := `
		SELECT id, vehicle_id, latitude, longitude, timestamp
		FROM vehicle_locations
		WHERE 1=1
	`
	queryArgs := pgx.NamedArgs{}

	if filter.Start != nil && filter.End != nil {
		query += " AND timestamp BETWEEN @start AND @end"
		queryArgs["start"] = filter.Start
		queryArgs["end"] = filter.End
	}
	if filter.Latest != nil && *filter.Latest {
		query += " ORDER BY timestamp DESC"
	}

	rows, err := r.db.Query(ctx, query, queryArgs)
	if err != nil {
		if err == pgx.ErrNoRows {
			return vehicleLocations, nil
		}
		return vehicleLocations, err
	}
	defer rows.Close()

	for rows.Next() {
		vehicleLocation := models.VehicleLocation{}
		err = rows.Scan(&vehicleLocation.ID, &vehicleLocation.VehicleID, &vehicleLocation.Latitude, &vehicleLocation.Longitude, &vehicleLocation.Timestamp)
		if err != nil {
			if err == pgx.ErrNoRows {
				return vehicleLocations, nil
			}
			return vehicleLocations, err
		}

		vehicleLocations = append(vehicleLocations, &vehicleLocation)
	}

	return vehicleLocations, nil
}
