FROM golang:1.20-alpine AS builder

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

###############################
FROM alpine

WORKDIR /artifacts

# copy
COPY db/migrations ./migrations
COPY --from=builder /usr/local/src/bin/ ./bin

VOLUME /artifacts

CMD ["./bin/gophkeeper"]