# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Set Go proxy to avoid checksum issues
ENV GOPROXY=https://proxy.golang.org,direct
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o expense-api ./cmd/server/main.go

# Stage 2: Run
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/expense-api .

EXPOSE 8085

CMD ["./expense-api"]
