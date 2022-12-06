FROM golang:latest
WORKDIR /app
COPY . .
RUN touch .env
RUN go mod download
RUN go build -o ./bin/api cmd/api/main.go

CMD ["./bin/api"]