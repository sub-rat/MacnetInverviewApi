FROM golang:latest
WORKDIR /app
COPY . .
RUN touch .env
RUN go mod download
RUN go build -o ./bin/macnet_api_assingment cmd/api/main.go

CMD ["./bin/macnet_api_assingment"]