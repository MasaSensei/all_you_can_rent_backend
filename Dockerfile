# ---- Build stage ----
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/api    ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/worker ./cmd/worker

# ---- API image ----
FROM alpine:3.20 AS api
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /bin/api /app/api
EXPOSE 8080
ENTRYPOINT ["/app/api"]

# ---- Worker image ----
FROM alpine:3.20 AS worker
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /bin/worker /app/worker
ENTRYPOINT ["/app/worker"]
