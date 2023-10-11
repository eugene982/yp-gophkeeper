BUILD_VERSION="v1.25.3"
BUILD_DATE=$(shell date +"%Y/%m/%d %H:%M")
BUILD_COMMIT=$(shell git rev-parse HEAD)

export PATH := $(PATH):$(shell go env GOPATH)/bin

BIN_PATH=./bin/gophkeeper

gofmt:
	gofmt -s -l . 	

tests:
	go test -race ./...

vet:
	go vet ./... 

staticcheck:
	staticcheck ./...

codecov:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...

golangci-lint:
	golangci-lint run ./...	

buildlint:
	go build -o=bin/staticlint cmd/staticlint/main.go 

runsrv:
	go run -ldflags \
		"-X main.buildVersion=$(BUILD_VERSION) -X 'main.buildDate=$(BUILD_DATE)' -X 'main.buildCommit=$(BUILD_COMMIT)' "\
		cmd/shortener/*.go

buildsrv: tests staticcheck vet
	go build -o $(BIN_PATH) \
		-ldflags \
		"-X main.buildVersion=$(BUILD_VERSION) -X 'main.buildDate=$(BUILD_DATE)' -X 'main.buildCommit="$(BUILD_COMMIT)"' "\
		cmd/shortener/*.go