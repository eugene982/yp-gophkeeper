BUILD_VERSION="v0.9.1"
BUILD_DATE=$(shell date +"%Y/%m/%d %H:%M")
BUILD_COMMIT=$(shell git rev-parse HEAD)

export PATH := $(PATH):$(shell go env GOPATH)/bin

DATABASE_DSN="postgres://test:test@localhost/gophkeeper_test"
VET_TOOL=./bin/statictest

gofmt:
	gofmt -s -l . 	

test:
	gofmt -s -l . \
	&& go test -race ./...

vet:
	go vet ./...

vettool:
	go vet -vettool=$$(which $(VET_TOOL)) ./... 

staticcheck:
	staticcheck ./...

codecov:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...

golangci-lint:
	golangci-lint run ./...	

# server
runsvr:
	go run -ldflags \
		"-X main.buildVersion=$(BUILD_VERSION) -X 'main.buildDate=$(BUILD_DATE)' -X 'main.buildCommit=$(BUILD_COMMIT)' "\
		cmd/gophkeeper/*.go -d $(DATABASE_DSN) 

buildsrv:
	go build -o ./bin/gophkeeper \
		-ldflags \
		"-X main.buildVersion=$(BUILD_VERSION) -X 'main.buildDate=$(BUILD_DATE)' -X 'main.buildCommit="$(BUILD_COMMIT)"' "\
		cmd/gophkeeper/*.go

# client
buildcli-lin:
		go build -o ./build/gk-client-linux \
		-ldflags \
		"-X main.buildVersion=$(BUILD_VERSION) -X 'main.buildDate=$(BUILD_DATE)' -X 'main.buildCommit="$(BUILD_COMMIT)"' "\
		cmd/grpcclient/*.go

buildcli-win:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
	go build -o ./build/gk-client-windows.exe \
		-ldflags \
		"-X main.buildVersion=$(BUILD_VERSION) -X 'main.buildDate=$(BUILD_DATE)' -X 'main.buildCommit="$(BUILD_COMMIT)"' "\
		cmd/grpcclient/main.go

buildcli-mac:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
	go build -o ./build/gk-client-windows \
		-ldflags \
		"-X main.buildVersion=$(BUILD_VERSION) -X 'main.buildDate=$(BUILD_DATE)' -X 'main.buildCommit="$(BUILD_COMMIT)"' "\
		cmd/grpcclient/main.go

runcli:
	go run cmd/grpcclient/*.go -l=debug

# gRPC
protoc: 
	protoc --go_out=gen/go --go_opt=paths=source_relative \
	--go-grpc_out=gen/go --go-grpc_opt=paths=source_relative \
	proto/v1/gophkeeper.proto

# migration
migrate-create:
	migrate create -ext sql -dir db/migrations -seq init schema

migrate-up:
	migrate -database $(DATABASE_DSN) -path db/migrations up

migrate-down:
	migrate -database $(DATABASE_DSN) -path db/migrations down

# Docker
docker-build:
	docker-compose up --build yp-gophkeeper

docker-run: docker-build
	docker-compose run --name yp-gophkeeper



