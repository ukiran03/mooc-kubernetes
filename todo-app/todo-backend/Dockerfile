FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-backend .

FROM alpine:3.21

# Install curl for debugging
# We do this as root before switching to the unprivileged user
RUN apk add --no-cache curl

RUN adduser -D appuser
USER appuser

WORKDIR /usr/src/app
COPY --from=builder /app/todo-backend .

CMD ["./todo-backend"]
