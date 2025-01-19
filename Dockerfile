# Use an official Go image as the base for building the app
FROM golang:1.20 AS build

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to the working directory
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o out main.go

# Use a lightweight image for running the app
FROM debian:bullseye-slim

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the build stage
COPY --from=build /app/out .

# Expose the port the app will run on
EXPOSE 3000

# Command to run the app
CMD ["./out"]
