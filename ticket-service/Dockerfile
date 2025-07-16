FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o ticket-service main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/ticket-service .
COPY env.example .env
EXPOSE 50051
CMD ["./ticket-service"] 