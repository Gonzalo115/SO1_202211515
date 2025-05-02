FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o grpc-server-rabbit .

FROM alpine:3.19
WORKDIR /root/
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/grpc-server-rabbit .
EXPOSE 50052
CMD ["./grpc-server-rabbit"]