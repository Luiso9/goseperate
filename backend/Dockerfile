FROM gocv:latest

# installing dep
RUN apt update && apt install -y python3-pip python3-venv && \
    python3 -m venv sklearn-env && \
    /sklearn-env/venv/bin/pip install numpy skicit-learn opencv-python pillow customtkinter pyperclip flask

# Set python
ENV PATH="/sklearn-env/venv/bin:$PATH"

# Set the working directory inside the container
WORKDIR /app

# Copy everything from the current directory into /app inside the container
COPY . .

# Build the Go application
RUN go build -o goseperate

# Command to run the binary
CMD ["./goseperate"]

