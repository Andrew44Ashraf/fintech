FROM golang:1.21

WORKDIR /app

# Ensure go.mod exists before copying files
RUN go mod init fintech-service

# Copy go.mod first to leverage Docker caching
COPY go.mod ./
RUN go mod tidy || true

# Copy the rest of the application files
COPY . .

# Build the application
RUN go mod init github.com/Andrew44Ashraf/fintech-service

EXPOSE 8080

CMD ["./fintech_service"]
