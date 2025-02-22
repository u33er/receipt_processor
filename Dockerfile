FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main cmd/server/main.go

FROM alpine:latest

WORKDIR /root/

ARG CONFIG_PATH

ENV CONFIG_PATH=${CONFIG_PATH}

COPY --from=builder /app/main .

COPY config /config

EXPOSE 8080

CMD ["./main"]