#!/usr/bin/env bash

export $(grep -v '^#' .env | xargs)
go run ../cmd/weather/main.go

