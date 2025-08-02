FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o atthehub ./cmd/ath


FROM alpine:latest AS runner

RUN apk --no-cache add ca-certificates

WORKDIR /app

RUN adduser -D -s /bin/sh appuser

COPY --from=builder /app/atthehub .

RUN chown appuser:appuser atthehub

USER appuser

CMD ["./atthehub"]