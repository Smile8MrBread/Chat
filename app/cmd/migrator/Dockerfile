FROM debian:bookworm

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends --assume-yes \
    sqlite3 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY app/migrations /app/migrations
COPY app/storage /app/storage

CMD /bin/sh -c "sqlite3 /app/storage/app.sqlite3 '.read /app/migrations/1_init.up.sql' && tail -f /dev/null"