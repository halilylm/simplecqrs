# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /bin/server ./app/product-api/main.go

## Deploy
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /bin/server /server
COPY --from=build /app/.env /.env

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/server"]