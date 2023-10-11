BUILD_VERSION="v0.0.1"
BUILD_DATE=$(shell date +"%Y/%m/%d %H:%M")
BUILD_COMMIT=$(shell git rev-parse HEAD)

export PATH := $(PATH):$(shell go env GOPATH)/bin

BIN_PATH=./bin/gophkeeper

gofmt:
	gofmt -s -l . 	

test:
	go test -race ./...

vet:
	go vet ./... 

staticcheck:
	staticcheck ./...

codecov:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...

golangci-lint:
	golangci-lint run ./...	

run:
	go run -ldflags \
		"-X main.buildVersion=$(BUILD_VERSION) -X 'main.buildDate=$(BUILD_DATE)' -X 'main.buildCommit=$(BUILD_COMMIT)' "\
		cmd/gophkeeper/*.go

buildsrv: tests staticcheck vet
	go build -o $(BIN_PATH) \
		-ldflags \
		"-X main.buildVersion=$(BUILD_VERSION) -X 'main.buildDate=$(BUILD_DATE)' -X 'main.buildCommit="$(BUILD_COMMIT)"' "\
		cmd/gophkeeper/*.go