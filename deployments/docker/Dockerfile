FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o go-family-finance ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/go-family-finance .
COPY --from=builder /app/deployments/config ./deployments/config

EXPOSE 8080

CMD ["./go-family-finance"] 