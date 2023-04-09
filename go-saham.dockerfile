FROM golang:1.17-alpine AS build
WORKDIR /app
COPY . /app

# Install dependencies
RUN apk update && apk add --no-cache git
RUN go mod tidy
RUN go get -d -v ./...
RUN go install -v ./...

# Build the binary
RUN go build -o main ./cmd/api

# Final stage
FROM alpine:latest

# Set environment variables for the database
ENV PORT=3000
ENV REDIS_PORT=localhost:6376
ENV REDIS_PASSWORD=eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81
ENV DBUser=go_saham
ENV DBPassword=go_saham
ENV DBName=go_saham
ENV DBHost=0.0.0.0
ENV DBPort=5436
ENV JWT_SECRET=notsecret

WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/main /app/main
COPY --from=build /app/.env /app/.env

EXPOSE 3000

# Start the server
CMD ["/app/main"]