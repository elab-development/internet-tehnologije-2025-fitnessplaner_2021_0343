FROM golang:1.22-alpine AS builder

WORKDIR /app

# Instalacija build alatki
RUN apk add --no-cache git

# Go module fajlovi
COPY go.mod ./
RUN go mod download

# Copy the rest of the backend source
COPY . .

# Osiguravanje go.sum je up to date & build
RUN go mod tidy && go build -o server ./main.go

FROM alpine:3.19

WORKDIR /app

# HTTPS konekcija
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server .
# Kopiraj OpenAPI specifikaciju za Swagger UI
COPY --from=builder /app/docs ./docs

# Default environment varijable 
ENV DB_USER=root \
    DB_PASSWORD=Vojislav123! \
    DB_HOST=db \
    DB_PORT=3306 \
    DB_NAME=app_db

EXPOSE 8080

CMD ["./server"]



