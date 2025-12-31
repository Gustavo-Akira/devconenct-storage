FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o devconnect-api ./cmd/api

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/devconnect-api .

EXPOSE 8080

CMD ["./devconnect-api"]
