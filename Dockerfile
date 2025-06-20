FROM golang:1.21  
WORKDIR /app

# Pertama copy mod/sum untuk caching
COPY git/go.mod .
COPY git/go.sum .
RUN go mod download

# Copy seluruh project
COPY git .

# Build dari folder cmd
RUN go build -o app ./cmd

EXPOSE 8080
CMD ["./app"]
