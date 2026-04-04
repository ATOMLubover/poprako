# ---- Build stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o poprako-s ./main.go

# ---- Runtime stage ----
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/poprako-s          ./poprako-s
COPY --from=builder /app/docs/              ./docs/

# Use the docker-specific config (listens on 0.0.0.0)
COPY docker/poprako-s-main/app_config.json  ./app_config.json

EXPOSE 8080

CMD ["./poprako-s"]
