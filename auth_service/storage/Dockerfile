FROM debian:bookworm

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends --assume-yes \
    sqlite3 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY auth_service/migrations /auth_service/migrations
COPY auth_service/storage/db /auth_service/storage

CMD /bin/sh -c "sqlite3 /auth_service/storage/auth_service.sqlite3 '.read /auth_service/migrations/1_init.up.sql' && tail -f /dev/null"