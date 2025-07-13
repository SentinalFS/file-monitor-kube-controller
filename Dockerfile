# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o file-monitor-controller .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/file-monitor-controller .

CMD ["./file-monitor-controller"]
