# Build stage
FROM golang:1.27-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/tiggerd ./cmd/tiggerd

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/tiggerd .
EXPOSE 8989
ENTRYPOINT ["./tiggerd"]