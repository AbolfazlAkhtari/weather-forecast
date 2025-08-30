.PHONY: weather-setup weather-run weather-build

weather-setup:
	@bash scripts/weather/setup.sh

weather-run:
	@bash scripts/weather/run.sh

weather-build:
	@bash scripts/weather/build.sh

migrate-create:
ifndef name
	$(error Please specify name, like: make migrate-create name=users)
endif
	goose -dir migrations create $(name) sql

migrate-up:
	@export $$(cat .env | xargs); \
	echo "Running migrations on $$DB_URL..."; \
	goose -dir migrations postgres "$$DB_URL" up