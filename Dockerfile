# Use the Go image for building
FROM golang:1.23 as build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files first to leverage caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the application files
COPY . .

# Build the application for the correct architecture (Linux, AMD64)
RUN GOOS=linux GOARCH=amd64 go build -o /app/out main.go

# Use a minimal image for running the app
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the built binary from the build stage
COPY --from=build /app/out .

# Ensure the binary is executable (this step is necessary if permissions need to be fixed)
RUN chmod +x /root/out

# Expose the application port
EXPOSE 3000

# Run the application
CMD ["/root/out"]
