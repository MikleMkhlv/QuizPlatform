CREATE TABLE room_players (
    room_id   UUID      NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    user_id   UUID      NOT NULL REFERENCES users(id),
    joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (room_id, user_id)
);