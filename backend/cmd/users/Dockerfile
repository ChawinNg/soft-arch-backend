# Use the official Golang image
FROM golang:1.23.1-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum from the shared root
COPY go.sum go.mod ./
RUN go mod download

# Copy the source code for the gRPC server
COPY . .

# Build the gRPC server
RUN go build -o users ./cmd/users/main.go

# Expose the gRPC port
EXPOSE 9000

# Run the gRPC server
CMD ["./users"]
