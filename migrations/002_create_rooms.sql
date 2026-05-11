CREATE TYPE room_status AS ENUM ('waiting', 'active', 'finished');

CREATE TABLE rooms (
    id          UUID        PRIMARY KEY,
    host_id     UUID        NOT NULL REFERENCES users(id), -- ← внешний ключ
    code        VARCHAR(10) NOT NULL UNIQUE,
    status      room_status NOT NULL DEFAULT 'waiting',
    max_players INTEGER     NOT NULL DEFAULT 10,
    created_at  TIMESTAMP   NOT NULL DEFAULT NOW()
);