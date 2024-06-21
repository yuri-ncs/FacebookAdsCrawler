# Stage 1: Build the Go binary
FROM golang:1.19-alpine as builder

# Install git (required for downloading Go modules)
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Log the contents of /app directory after downloading dependencies
RUN ls -la /app

# Copy the rest of the application source code
COPY . .

# Log the contents of /app directory after copying source code
RUN ls -la /app

# Ensure main.go exists and build the Go binary
RUN go build -o /app/ads-api .

# Log the contents of /app directory after building the binary
RUN ls -la /app

# Stage 2: Create a smaller image with the built binary
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Install certificates for making HTTPS requests if needed
RUN apk --no-cache add ca-certificates

# Copy the built binary from the builder stage
COPY --from=builder /app/ads-api .

# Log the contents of /root directory after copying the binary
RUN ls -la /root/

# Expose the port the service will run on
EXPOSE 8000

# Command to run the executable
CMD ["./ads-api"]