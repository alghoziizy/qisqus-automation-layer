# Dockerfile
FROM golang:1.23

# Set working directory
WORKDIR /app

# Copy go mod dan sum dulu (untuk cache build layer dependency)
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy semua source code ke dalam container
COPY . .

# Build app
RUN go build -o app .

# Expose port (default Gin = 8080)
EXPOSE 8080

# Jalankan binary
CMD ["./app"]
