# Build stage
FROM golang:1.21-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

EXPOSE 3000
CMD ["/app/main"]