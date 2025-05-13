# Source https://habr.com/ru/companies/otus/articles/660301/

FROM golang:alpine@sha256:ef18ee7117463ac1055f5a370ed18b8750f01589f13ea0b48642f5792b234044 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED=0 \
    GOOS=linux

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .

RUN mkdir -p /build/app && go build -ldflags="-s -w" -o /build/app/notification cmd/notification/main.go

# Final image
FROM alpine:latest@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c

RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /build/app/notification /app/notification

CMD ["/app/notification"]
