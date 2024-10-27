# Virtual currency service

Provides a thin additional layer above payment systems. The layer stores currency balances of users. The service functionality may be extended by modifying the request pipeline handler.

## Usage

* Specify configuration in a *.env* file at the root of the project. The config should include the following fields:
  - ```SERVE_ENDPOINT```
  - ```SERVE_PORT```
  - ```DB_HOST```
  - ```DB_PORT```
  - ```DB_USER```
  - ```DB_PASSWORD```
  - ```DB_NAME```
* Note, that for now only the postgres database is supported
* Build the service:
  - Docker: ```sudo docker build -t mse-back:1.0 .  ```
  - Executable: ```go build -o srvr ./cmd/main.g```
* Run the service with:
  - Docker: ```sudo docker run -d -p 8787:8787 --network <db_network>  mse-back:1.0```
* Setup requests handling logic: each user balance modification request is handled by this service and external services in a specific order. This order is declared in an implementation of the *Pipeline* pattern in the virtual currency service. You can modify the pipeline using the api, described in the currency-project-api/openapi-back.yaml. There are currenly 3 types of pipelines:
  - *Statistics* service, being called first to observe the state of the system,
  - *Authentication* services, being called second to control the access,
  - *Payment* service, being called at the end to perform the payment.

## Development

### Technologies

- The [Go](https://go.dev/tour/welcome/1) language
- The [gin](https://github.com/gin-gonic/gin) router
- The [Docker](https://docs.docker.com/) containerization system
- The [Postgres](https://www.postgresql.org/docs/) database
- HTTP Protocol

### Modifying

1. Create a *.env.development* file at the root of the project containing the *.env* file fields
2. ```go run ./cmd/main.go --dev```

## TODO

1. Add support for gRPC
2. Add support for databases besides Posgres
