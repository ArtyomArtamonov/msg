CREATE TABLE rooms (
    id UUID PRIMARY KEY,
    name VARCHAR(30) NOT NULL DEFAULT 'DEFAULT_ROOM_NAME',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    dialog_room BOOLEAN NOT NULL
);

CREATE TABLE user_in_room (
    room_id UUID NOT NULL,
    user_id UUID NOT NULL,
    PRIMARY KEY(room_id, user_id),
    CONSTRAINT fk_room_id
        FOREIGN KEY(room_id)
            REFERENCES rooms(id)
            ON DELETE CASCADE,
    CONSTRAINT fk_user_id
        FOREIGN KEY(user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
);

CREATE TABLE messages (
    id UUID NOT NULL PRIMARY KEY,
    room_id UUID NOT NULL,
    user_id UUID NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    CONSTRAINT fk_room_id
        FOREIGN KEY(room_id)
            REFERENCES rooms(id)
            ON DELETE CASCADE,
    CONSTRAINT fk_user_id
        FOREIGN KEY(user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
);
