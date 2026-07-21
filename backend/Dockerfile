# --- Development container ---
FROM golang:1.26.5-alpine3.24

WORKDIR /app

RUN apk add --no-cache git build-base

# Air
RUN go install github.com/air-verse/air@latest

# dependencies
COPY go.mod go.sum* ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]