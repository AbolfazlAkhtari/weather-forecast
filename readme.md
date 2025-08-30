# Weather App
This project contains a Go-based application with Dockerized infrastructure.


## Prerequisites
1. Copy .env.example and make .env : `cp .env.example .env`
2. Docker & Docker Compose
3. Go (Only required for local development)


## Available Commands
We use a `Makefile` to simplify common tasks. The scripts are located under `scripts/weather/`.


## Step-by-Step Guide to Run the app
### 1. Setup Infrastructure

Sets up the required infrastructure (Postgres DB for now) using Docker Compose:

```sh   
make weather-setup
```


### 2. Run Application
- Run With Go
  - Recommended for developers: no build time, LOCAL only
  - running the app with `go run`. needs go to be installed on local
  - if you wish to build the application, jump to the next step
  - ```sh
    make weather-run
    ```
- Build With Docker
  - Build the application and run it in a docker cotainer
  - ```sh
    make weather-build
    ```


Now the application is up and running on `localhost:{port}` 

-----> **{port} = `WEATHER_PORT` set in env**


### 3. Migrations
Migrations are handled via [goose library](https://github.com/pressly/goose)

* Migrations Run on each build of the application so if you just need to run the app without 
any need to change the database schema, you can totally skip below steps,

to install it, run:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Create a new migration
```bash
make migrate-create name=weathers
```

### Run Migrations
Migrations automatically run on each build of the application.
```bash
make migrate-up
```