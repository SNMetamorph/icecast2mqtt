# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder

RUN go install github.com/go-task/task/v3/cmd/task@latest

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN task build

FROM alpine:3.20

RUN apk add --no-cache ca-certificates && \
    addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY config.schema.json ./
COPY config.json.example ./
COPY --from=builder /src/bin/icecast2mqtt /usr/local/bin/icecast2mqtt

USER app

CMD ["icecast2mqtt"]