# Stage 1: Build
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY .env .

COPY . .


RUN go mod download




RUN go build -o main ./cmd/main.go


FROM alpine:latest


COPY --from=builder /app/main /app/main
COPY --from=builder /app/.env /app/.env


WORKDIR /app

EXPOSE 8080

CMD ["./main"]
