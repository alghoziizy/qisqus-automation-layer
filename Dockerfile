FROM golang:1.21

WORKDIR /app

# Copy mod/sum pertama untuk caching
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy semua file
COPY . .

# Build aplikasi
RUN go build -o app ./cmd

EXPOSE 8080
CMD ["./app"]
