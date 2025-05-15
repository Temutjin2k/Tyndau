FROM golang:alpine AS builder

WORKDIR /build

ENV CGO_ENABLED=0 \
    GOOS=linux

RUN apk update --no-cache && apk add --no-cache tzdata

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o music-service ./cmd/music/main.go

FROM alpine:latest

RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /build/music-service .

CMD ["./music-service"]
