FROM golang:1.16.9-alpine3.14 AS builder

RUN apk add gcc libc-dev ca-certificates linux-headers git

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN swag init -g server.go

RUN go build -o server *.go

FROM alpine:3.12.0

RUN apk add --no-cache ca-certificates curl

WORKDIR /app

COPY --from=builder /config/config.json ./config/config.json
COPY --from=builder /app/server .

COPY --from=builder /app/docs/ /app/docs/

RUN chmod +x ./server
CMD ["./server"]
