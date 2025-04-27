# Event-Driven Architecture Makefile

.PHONY: build-api build-publisher build-worker run-api run-publisher run-worker test

build-api:
	go build -o bin/api cmd/api/main.go

build-publisher:
	go build -o bin/publisher cmd/publisher/main.go

build-worker:
	go build -o bin/worker cmd/worker/main.go

run-api: build-api
	./bin/api

run-publisher: build-publisher
	./bin/publisher

run-worker: build-worker
	./bin/worker

test:
	go test ./...

clean:
	rm -rf bin/
