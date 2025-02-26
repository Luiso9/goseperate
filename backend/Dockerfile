# Use a lightweight base image with Go
FROM gocv:latest

# Set the working directory inside the container
WORKDIR /app

# Copy everything from the current directory into /app inside the container
COPY . .

# Build the Go application
RUN go build -o goseperate

# Command to run the binary
CMD ["./goseperate"]

