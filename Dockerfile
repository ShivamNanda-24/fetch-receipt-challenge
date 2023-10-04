# Use an official Go runtime as a parent image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the rest of the application code into the container
COPY . .

# Build the Go application
RUN go build -o main

# Expose a port if your application listens on a specific port
EXPOSE 8080

# Command to run the application
CMD ["./main"]
