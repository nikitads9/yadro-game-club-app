BIN_APP := "./bin/events"
.PHONY: build
build: build-server build-client

build-server:
	go build -v -ldflags "-w -s" -o $(BIN_APP) ./cmd/events/events.go

.PHONY: run
run: env docker-compose-up
docker-compose-up:
	docker-compose up -d

.PHONY: stop
stop:
	docker compose stop