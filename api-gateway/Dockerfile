# Source https://habr.com/ru/companies/otus/articles/660301/

FROM golang:alpine@sha256:7772cb5322baa875edd74705556d08f0eeca7b9c4b5367754ce3f2f00041ccee AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED=0 \
    GOOS=linux

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .

RUN mkdir -p /build/app && go build -ldflags="-s -w" -o /build/app/gateway cmd/gateway/main.go

# Final image
FROM alpine:latest@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c

RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /build/app/gateway /app/gateway

CMD ["/app/gateway"]
