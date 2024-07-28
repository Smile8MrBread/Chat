FROM golang:bookworm AS builder

WORKDIR /chat_app

COPY client/go.mod client/go.sum ./
RUN go mod download

COPY client .

RUN go build -o /client ./cmd/main.go

FROM debian:bookworm

COPY --from=builder /client /usr/bin/client
COPY --from=builder /chat_app/cmd/frontend /client/cmd/frontend
COPY --from=builder /chat_app/cmd/tmp /client/cmd/tmp
COPY --from=builder /chat_app/config/local.yaml /config/config.yaml

CMD ["/usr/bin/client"]
