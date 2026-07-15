# --- Development container ---
FROM golang:1.26.5-alpine3.24

WORKDIR /app

# Server variables
ENV SERVER_ADDRESS=0.0.0.0:8080
# Timeout in seconds
ENV SERVER_TIMEOUT=4s
ENV SERVER_IDLE_TIMEOUT=60s
ENV SERVER_SHUTDOWN_TIMEOUT=10s

# DB variables
ENV DB_NAME=postgres
ENV DB_USER=dbuser
ENV DB_PASS=my_db_password

RUN apk add --no-cache git build-base

# Air
RUN go install github.com/air-verse/air@latest

# dependencies
COPY go.mod go.sum* ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]