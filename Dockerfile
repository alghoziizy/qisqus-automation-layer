FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN cd cmd && go build -o ../out

FROM alpine
WORKDIR /app
COPY --from=builder /app/out .
COPY .env .
CMD ["./out"]
