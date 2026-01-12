FROM golang:1.25.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-app .


FROM alpine:latest
WORKDIR /usr/src/app
COPY --from=builder /app/todo-app .
CMD ["./todo-app"]
