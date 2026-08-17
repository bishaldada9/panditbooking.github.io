FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /out/server ./cmd/main.go \
    && CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /out/seed ./scripts/seed.go

FROM alpine:3.19 AS runtime

RUN apk add --no-cache ca-certificates tzdata wget

WORKDIR /app

COPY --from=builder /out/server .

EXPOSE 8080

CMD ["./server"]

FROM runtime AS seed

COPY --from=builder /out/seed .

CMD ["./seed"]
