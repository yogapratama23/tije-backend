
# Backend Technical Test

### How to Start
- Clone this repository
- Run docker compose
```
    // Inside the repository

    $ docker compose up --build -d

    // Go API will be running in port 8010
    // you can customize the port by editing docker-compose.yml
```

#### Services list
- Golang API using Gin
- RabbitMQ Receiver
- PostgresQL
- RabbitMQ
- Eclipse Mosquitto

#### Golang API
- Simulate Driving
``` 
GET {{host}}/simulate-driving?vehicle_id={vehicleId} 
```
This endpoint will trigger driving simulation that will insert driving coordinates to the database for every 2 seconds, everytime the coordinates is sent, it will also calculate the distance to each checkpoint, and if one of the checkpoint distance is below 50 meters, it will trigger geofence_alert to RabbitMQ and RabbitMQ Receiver will log the alert

- Latest Location
```
GET {host}/vehicles/{vehicleId}/location
```
This endpoint shows the latest location of the selected vehicle id

- Location History
```
GET {{host}}/vehicles/{vehicleId}/history?start=1750960862&end=1750960864
```
This endpoint shows history of vehicle locations



