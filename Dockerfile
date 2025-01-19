# Use the Go image for building
FROM golang:1.23 as build

# Set the working directory inside the container
WORKDIR /app

# Copy all files into the container
COPY . .

# Download dependencies
RUN go mod tidy

# Build the application
RUN go build -o out main.go

# Use a minimal image for running the app
FROM alpine:latest
WORKDIR /root/
COPY --from=build /app/out .

# Expose the application port
EXPOSE 3000

# Run the application
CMD ["./out"]
