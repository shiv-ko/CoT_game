# Use an official Go runtime as a parent image
FROM golang:1.24-alpine

# Set the working directory in the container
WORKDIR /app

# Copy Go module and dependency files
COPY go.mod go.sum ./

# Download and install Go dependencies
RUN go mod download

# Copy the rest of the application code to the container
COPY . .

CMD ["go", "run", "main.go"]
