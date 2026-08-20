DROP TABLE IF EXISTS avatars;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    user_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    vrc_name VARCHAR(16) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);
CREATE TABLE avatars (
    avatar_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    avatar_name TEXT NOT NULL,
    creator TEXT NOT NULL,
    shop_url TEXT,
    memo TEXT
);