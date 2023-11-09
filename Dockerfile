FROM golang:1.20-alpine

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc musl-dev postgresql-client

# download dependensies
COPY ["go.mod", "go.sum", "./"]
RUN go mod download

# build
COPY cmd ./cmd
COPY gen ./gen
COPY internal ./internal

RUN mkdir ./bin
RUN go build -o ./bin/gophkeeper cmd/gophkeeper/main.go

#copy migrations
COPY db/migrations ./migrations

# start
RUN chmod +x ./bin/gophkeeper
CMD ["./bin/gophkeeper"]
