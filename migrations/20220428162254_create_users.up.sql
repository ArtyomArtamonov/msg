CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(25) NOT NULL,
    password_hash VARCHAR(64) NOT NULL,
    role VARCHAR(25) DEFAULT 'user' NOT NULL
);
