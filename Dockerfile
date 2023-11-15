FROM golang:1.20.11-alpine3.18 AS builder

RUN apk add gcc libc-dev ca-certificates linux-headers git

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o server cmd/app/*.go

FROM alpine:3.12.0

RUN apk add --no-cache ca-certificates curl

WORKDIR /app

COPY --from=builder /app/config/config.json ./config/config.json
COPY --from=builder /app/server .

RUN chmod +x ./server
CMD ["./server"]
