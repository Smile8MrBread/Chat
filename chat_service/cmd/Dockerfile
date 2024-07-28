FROM golang:bookworm AS builder

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends --assume-yes \
    build-essential \
    libsqlite3-dev

WORKDIR /chat_app

COPY chat_service/go.mod chat_service/go.sum ./
RUN go mod download

COPY chat_service .

RUN go build -o /chat_service ./cmd/main.go

FROM debian:bookworm

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends --assume-yes \
    sqlite3 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /chat_service /usr/bin/chat_service
COPY --from=builder /chat_app/config/local.yaml /config/config.yaml

EXPOSE 32200

CMD ["/usr/bin/chat_service"]
