FROM golang:1.21.0 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Build the Go app
RUN go build -o /server

# Start a new stage from scratch
FROM gcr.io/distroless/base-debian10

# Set the Current Working Directory inside the container
WORKDIR /

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /server /server

# Expose port 8080 to the outside world
EXPOSE 8080

# Use an unprivileged user
USER nonroot:nonroot

# Command to run the executable
ENTRYPOINT [ "/server" ]
