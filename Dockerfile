# Build stage
FROM golang:1.25.3 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make install

# Final stage
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/bin/webserver .
EXPOSE 8080
ENV LANG=C.UTF-8
ENV GIN_MODE=release
CMD ["./webserver"]
