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

# Build the application statically (ensure no dependencies on libraries)
RUN CGO_ENABLED=0 GOOS=linux go build -o out main.go

# Runtime Stage
FROM alpine:latest

# Install dependencies needed to run the Go binary in Alpine
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container for runtime
WORKDIR /root/

# Copy the built Go binary from the build stage
COPY --from=build /app/out .

# Expose the port your application will listen on
EXPOSE 3000

# Command to run your application
CMD ["/root/out"]
