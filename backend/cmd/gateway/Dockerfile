# Use the official Golang image
FROM golang:1.23.1-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum from the shared root
COPY go.sum go.mod ./
RUN go mod download

# Copy the source code for the API Gateway
COPY . .

# Build the API gateway
RUN go build -o gateway ./cmd/gateway/main.go

# Expose the HTTP port
EXPOSE 8080

# Run the API Gateway
CMD ["./gateway"]
