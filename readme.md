# Weather App

This project contains a Go-based application with Dockerized infrastructure.

## Prerequisites

- Go (for local development / `make weather-run`)
- Docker & Docker Compose

## Available Commands

We use a `Makefile` to simplify common tasks. The scripts are located under `scripts/weather/`.

### 1. Setup Infrastructure (Required Step)

Sets up the required infrastructure (Postgres DB, etc.) using Docker Compose:

```sh   
make weather-setup
```

### 2. Run Application (Local Go)
*(if you wish to build the application, jump to the next step)*

Runs the application locally with your Go installation:

```sh
make weather-run
```

### 3. Build Application (Docker)
Builds the application image with Docker:

```sh
make weather-build
```
