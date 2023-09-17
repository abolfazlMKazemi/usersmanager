# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY . .

# Build the Go application
#RUN go build -o myapp ./cmd
RUN go build -o chargecode ./cmd
# Expose the port your Go application is listening on
EXPOSE 4238

# Run the Go application
CMD ["./chargecode"]
