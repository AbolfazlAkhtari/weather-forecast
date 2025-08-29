.PHONY: weather-setup weather-run weather-build

weather-setup:
	@bash scripts/weather/setup.sh

weather-run:
	@bash scripts/weather/run.sh

weather-build:
	@bash scripts/weather/build.sh
