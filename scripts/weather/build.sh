#!/usr/bin/env bash

export $(grep -v '^#' .env | xargs)
docker compose -f ../docker/weather/compose.yml up -d

