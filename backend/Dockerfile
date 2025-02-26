FROM golang:1.24-bookworm as BASE

# Set working directory
WORKDIR /app

# copy go.mod and go.sum
COPY go.mod go.sum ./

RUN go mod download

COPY . .

# build
RUN go build -o goseperate

# Expose the port
EXPOSE 8080

# Command to run the app
CMD ["/app/goseperate"]
