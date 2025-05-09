# Builder stage
FROM golang:1.23.5 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Final image
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/main .
COPY .env .  

RUN chmod +x /root/main

EXPOSE 8080
CMD ["/root/main"]

