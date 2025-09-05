# Сборка приложения
FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/quiz-please-scheduler .

# Финальный образ
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/quiz-please-scheduler .
# Copy production config to /configs inside the image
COPY configs/config.prod.yaml /configs/config.prod.yaml
CMD ["./quiz-please-scheduler"]