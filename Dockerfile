# Build Stage
FROM golang:1.23 as build

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the application
COPY . .

# Build the application
RUN go build -o out main.go

# Runtime Stage
FROM alpine:latest

# Set the working directory inside the container for runtime
WORKDIR /root/

# Copy the built Go binary from the build stage
COPY --from=build /app/out .

# Expose the port your application will listen on
EXPOSE 3000

# Command to run your application
CMD ["/root/out"]
