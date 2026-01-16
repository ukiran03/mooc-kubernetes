FROM golang:1.25.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o log-output .

FROM alpine:latest
WORKDIR /usr/src/app
COPY --from=builder /app/log-output .
CMD ["./log-output"]
