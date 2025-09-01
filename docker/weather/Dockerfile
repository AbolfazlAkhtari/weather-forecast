FROM golang:1.24

RUN apt-get update && apt-get install -y \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /builds/weather cmd/weather/main.go

EXPOSE 8000

# Run
CMD ["/builds/weather"]
