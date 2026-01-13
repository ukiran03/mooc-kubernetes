FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o todo-app .


FROM alpine:3.21

RUN adduser -D appuser
USER appuser

WORKDIR /usr/src/app

COPY --from=builder /app/ui ./ui
COPY --from=builder /app/todo-app .

CMD ["./todo-app"]
