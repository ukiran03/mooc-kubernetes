FROM golang:1.25.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o get-data .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/get-data .
CMD ["./get-data"]
