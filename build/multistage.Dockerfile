# syntax=docker/dockerfile:1

## Build
FROM golang:1.18-buster AS build

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./config ./config
COPY ./internal ./internal

WORKDIR /app/cmd/server

RUN go build -o /balance-service

## Deploy
FROM debian

WORKDIR /

COPY --from=build /balance-service /balance-service
COPY ./config/docker-config.yml /config/docker-config.yml

#USER nonroot:nonroot

ENTRYPOINT ["/balance-service"]