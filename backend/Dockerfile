FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build tools
RUN apk add --no-cache git

# Go module files
COPY go.mod ./
RUN go mod download

# Copy the rest of the backend source
COPY . .

# Ensure go.sum is up to date and build
RUN go mod tidy && go build -o server ./main.go

FROM alpine:3.19

WORKDIR /app

# Certificates for HTTPS connections if needed
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server .
# Kopiraj OpenAPI specifikaciju za Swagger UI
COPY --from=builder /app/docs ./docs

# Default environment variables (can be overridden in docker-compose)
ENV DB_USER=root \
    DB_PASSWORD=Vojislav123! \
    DB_HOST=db \
    DB_PORT=3306 \
    DB_NAME=app_db

EXPOSE 8080

CMD ["./server"]



