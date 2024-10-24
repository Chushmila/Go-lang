# Step 1: Use the official Golang image as the base image for building the app
FROM golang:1.20 AS builder

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy go.mod and go.sum files into the container
COPY go.mod go.sum ./

# Step 4: Download the dependencies
RUN go mod download

# Step 5: Copy the source code into the container
COPY . .

# Step 6: Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o helloWorld .

# Step 7: Create a new minimal image to run the binary (using Alpine Linux)
FROM alpine:latest

# Step 8: Set the working directory in the final image
WORKDIR /root/

# Step 9: Copy the built binary from the builder stage
COPY --from=builder /app/helloWorld .

# Step 10: Expose the port the app will listen on (if applicable)
EXPOSE 8080

# Step 11: Define the command to run when the container starts
CMD ["./helloWorld"]
 
