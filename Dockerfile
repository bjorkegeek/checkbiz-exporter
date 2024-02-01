# First stage: build the application
FROM golang:1.18 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and download dependencies
COPY go.* ./
RUN go mod download

# Copy the rest of the application's source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Second stage: create the runtime image
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Run the binary
CMD ["./main"]
