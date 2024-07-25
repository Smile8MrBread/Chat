FROM golang:bookworm AS builder

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends --assume-yes \
    build-essential \
    libsqlite3-dev

WORKDIR /chat_app

COPY go.mod go.sum ./
RUN go mod download

COPY ../cmd/services/auth .

RUN go build -o /auth_service ./app/cmd/services/auth_service/main.go

FROM debian:bookworm

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends --assume-yes \
    sqlite3 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /auth /usr/bin/auth
COPY --from=builder /chat_app/app/config/local.yaml /config/config.yaml

EXPOSE 32100

CMD ["/usr/bin/auth"]
