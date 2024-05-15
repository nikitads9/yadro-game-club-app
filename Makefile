include .env
BIN_APP := "./bin/events"

.PHONY: build
build: 
	CGO_ENABLED=0 GOOS=linux go build -v -ldflags "-w -s" -o $(BIN_APP) ./cmd/events/events.go

.PHONY: build-win
build-win: 
	CGO_ENABLED=0 GOOS=windows go build -v -ldflags "-w -s" -o $(BIN_APP).exe ./cmd/events/events.go

.PHONY: run
run: docker-build docker-run
docker-build:
	docker build --no-cache -f ./deploy/Dockerfile . --tag nikitads9/yadro-game-club:app
docker-run:
	set -o allexport && source ./.env && set +o allexport
	docker run -d -e DATA_PATH=${DATA_PATH} -v ${PWD}/testdata/:/testdata/ --name app nikitads9/yadro-game-club:app

.PHONY: wipe
wipe: docker-remove docker-delete
docker-remove:
	docker container rm app
docker-delete:
	docker image rm nikitads9/yadro-game-club:app