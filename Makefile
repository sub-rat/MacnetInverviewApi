run: build
	./bin/macnet_api_assingment

build:
	go build -o ./bin/macnet_api_assingment cmd/api/main.go

clean:
	go clean

test:
	go test ./...