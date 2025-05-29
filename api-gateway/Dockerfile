FROM golang:1.24.3

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY . .

# Download the dependencies
RUN go mod tidy

# Build the Go application
RUN go build -o main .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"]