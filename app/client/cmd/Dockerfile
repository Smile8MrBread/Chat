FROM golang:bookworm AS builder

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends --assume-yes \
    build-essential \
    libsqlite3-dev

WORKDIR /chat_app

COPY go.mod go.sum ./
RUN go mod download

COPY ../../cmd/client .

RUN go build -o /client ./app/cmd/client/main.go

FROM debian:bookworm

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends --assume-yes \
    sqlite3 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /chat_app/app/cmd/client/frontend /app/cmd/client/frontend
COPY --from=builder /chat_app/app/cmd/client/tmp /app/cmd/client/tmp
COPY --from=builder /client /usr/bin/client
COPY --from=builder /chat_app/app/config/local.yaml /config/config.yaml

CMD ["/usr/bin/client"]
