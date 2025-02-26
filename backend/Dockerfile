FROM debian:bookworm-slim

# Set working directory
WORKDIR /app

# Copy the manually built binary
COPY build .

# Expose the port
EXPOSE 8080

# Command to run the app
CMD ["./build"]
