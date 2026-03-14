# --- build stage ---
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o garden-app .

# --- runtime stage ---
FROM alpine:3.20

RUN addgroup -S garden && adduser -S -G garden garden
RUN mkdir -p /data && chown garden:garden /data

WORKDIR /app
COPY --from=builder /app/garden-app .

USER garden
EXPOSE 8080

# DB stored on mounted volume at /data/garden.db
CMD ["./garden-app", "serve", "--db", "/data/garden.db", "--port", "8080"]
