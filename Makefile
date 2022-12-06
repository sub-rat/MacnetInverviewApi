run: build
	./bin/api

build:
	go build -o ./bin/api cmd/api/main.go

clean:
	go clean

test:
	go test ./...