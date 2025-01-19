# Stage 1: Build
FROM golang:1.23 as builder
WORKDIR /app
COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download
COPY ./backend .
RUN go build -o estimator cmd/app/main.go

# Stage 2: Production
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/bin/estimator .
EXPOSE 8080
CMD ["./estimator"]

# Stage 3: Development
FROM golang:1.23 as dev
WORKDIR /app

COPY ./backend/.air.toml ./

# Install air for hot reloading
RUN go install github.com/air-verse/air@latest

# Copy source code
COPY ./backend .

# Expose application port
EXPOSE 8181

# Default command for development
CMD ["air", "-c", ".air.toml"]
