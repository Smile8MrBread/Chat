CREATE TABLE IF NOT EXISTS Chats
(
    id         INTEGER PRIMARY KEY,
    name       TEXT    NOT NULL UNIQUE,
    creator_id INTEGER NOT NULL REFERENCES Users (id),
    preview    TEXT DEFAULT NULL,
    avatar     TEXT DEFAULT 'default.jpg'
);

CREATE TABLE IF NOT EXISTS Messages
(
    id      INTEGER PRIMARY KEY,
    text    TEXT     NOT NULL,
--     image
    date    DATETIME NOT NULL,
    user_from INTEGER  NOT NULL REFERENCES Users (id),
    user_to INTEGER  NOT NULL REFERENCES Users (id)
);

CREATE TABLE IF NOT EXISTS Contacts
(
    user_id    NOT NULL REFERENCES Users (id),
    contact_id NOT NULL REFERENCES Users (id),
    is_messaged INT DEFAULT NULL
);