# build
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o deploy/MyStonksDao .

# run
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/deploy/MyStonksDao ./deploy/MyStonksDao
COPY --from=builder /app/config ./config

RUN chmod +x ./deploy/MyStonksDao

EXPOSE 8000
CMD ["./deploy/MyStonksDao", "start", "--config=config/config.yaml"]
