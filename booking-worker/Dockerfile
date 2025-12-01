FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o booking-worker main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/booking-worker .
COPY env.example .env
EXPOSE 50059
EXPOSE 9091
CMD ["./booking-worker"]

