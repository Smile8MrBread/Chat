FROM golang:bookworm AS builder

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends --assume-yes \
    build-essential \
    libsqlite3-dev

WORKDIR /chat_app

COPY auth_service/go.mod auth_service/go.sum ./
RUN go mod download

COPY auth_service .

RUN go build -o /auth_service ./cmd/main.go

FROM debian:bookworm

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends --assume-yes \
    sqlite3 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /auth_service /usr/bin/auth_service
COPY --from=builder /chat_app/config/local.yaml /config/config.yaml

EXPOSE 32100

CMD ["/usr/bin/auth_service"]
