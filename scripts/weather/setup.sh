#!/usr/bin/env bash

export $(grep -v '^#' .env | xargs)
docker compose -f docker/weather/infra.yml up -d

