FROM golang:bookworm AS builder

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends --assume-yes \
    build-essential \
    libsqlite3-dev

WORKDIR /chat_app

COPY go.mod go.sum ./
RUN go mod download

COPY ../cmd/services/chat .

RUN go build -o /chat_service ./app/cmd/services/chat_service/main.go

FROM debian:bookworm

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends --assume-yes \
    libsqlite3-0 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /chat /usr/bin/chat
COPY --from=builder /chat_app/app/config/local.yaml /config/config.yaml

EXPOSE 32200

CMD ["/usr/bin/chat"]
