# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Salin go.mod dan go.sum terlebih dahulu untuk memanfaatkan cache Docker
COPY go.mod go.sum ./
RUN go mod download

# Salin kode sumber
COPY . .

COPY internal/infrastructure/broadcast/donora-f67f2-5c889d5acd0a.json internal/infrastructure/broadcast/donora-f67f2-5c889d5acd0a.json

# Build aplikasi
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates && \
    update-ca-certificates

WORKDIR /root/

# Salin binary dari build stage
COPY --from=builder /app/main .

COPY --from=builder /app/internal/infrastructure/broadcast/donora-f67f2-5c889d5acd0a.json internal/infrastructure/broadcast/donora-f67f2-5c889d5acd0a.json

# Expose port (sesuaikan dengan kebutuhan aplikasi Anda)
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]