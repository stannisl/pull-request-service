FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service ./cmd/service

FROM scratch

WORKDIR /app

COPY --from=builder /app/service .

EXPOSE 8080

ENTRYPOINT ["./service"]