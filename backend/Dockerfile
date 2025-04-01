FROM python:3.12-slim AS base

# Install system dependencies and Go
RUN apt-get update && \
    apt-get install -y libglib2.0-0 libgl1-mesa-glx wget && \
    rm -rf /var/lib/apt/lists/*

# Install Go manually (replace 1.23 with the desired version if needed)
ENV GOLANG_VERSION=1.23.6
RUN wget https://go.dev/dl/go1.23.6.linux-amd64.tar.gz && \
    rm -rf /usr/local/go && \
    tar -C /usr/local -xzf go1.23.6.linux-amd64.tar.gz && \
    rm go1.23.6.linux-amd64.tar.gz


# Set Go environment variables
ENV PATH="/usr/local/go/bin:$PATH"
ENV GOPATH=/go

# Install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Set up Go workspace
WORKDIR /go/src/app

# Copy Go module files (go.mod, go.sum)
COPY go.mod go.sum ./

# Fetch Go dependencies
RUN go mod tidy

# Copy the Go source code (handlers, services, etc.)
COPY . /go/src/app

# Build Go project
RUN go build -o /app/go-backend main.go

# Copy Python files to /app (Python files should reside here for proper execution)
WORKDIR /app
COPY scripts /app/scripts

# Set the working directory for running the Python app
WORKDIR /app

# Run Python script and Go binary concurrently
CMD ["sh", "-c", "python /app/scripts/preview.py & /app/go-backend"]