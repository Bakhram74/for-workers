FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/api/main.go

# Run stage
FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
COPY config.yaml .
COPY db/migration ./db/migration

CMD [ "/app/main" ]