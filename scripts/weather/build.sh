#!/usr/bin/env bash

export $(grep -v '^#' .env | xargs)
docker network create weather_network
docker compose -f docker/weather/compose.yml up --build -d

