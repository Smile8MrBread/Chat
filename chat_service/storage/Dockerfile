FROM debian:bookworm

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends --assume-yes \
    sqlite3 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY chat_service/migrations /chat_service/migrations
COPY chat_service/storage/db /chat_service/storage

CMD /bin/sh -c "sqlite3 /chat_service/storage/chat_service.sqlite3 '.read /chat_service/migrations/1_init.up.sql' && tail -f /dev/null"