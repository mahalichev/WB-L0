# NATS-streaming service
Generation, receipt, validation and caching of orders using nats-streaming.

## Technologies
* Golang 1.21.6
* PostgreSQL 16.1 + GORM
* NATS-streaming 0.25.6 + STAN
* go-playground/validator
* joho/godotenv

## How to use
To start the server enter the following commands into the console:
```console
docker compose up -d --build
go run ./service/main.go
go run ./publisher/publisher.go
```
and then go to the address __localhost:3000__ (or __127.0.0.1:3000__).

To start load testing, change the value of the VEGETA_ORDERUID field in the .env file to the order_uid existing in the database and enter the following commands into console:
```console
go run ./service/main.go
go run ./load_testing/vegeta.go
```
GET requests will be sent to the server at a frequency of 1000 per second for 3 seconds.

## API methods
### GET /orders?id=\<id\>
Information about the order with order_uid == \<id\> in JSON format.

### GET /static/orders?id=\<id\>
HTML page with information about the order with order_uid == \<id\> and an order search field.
