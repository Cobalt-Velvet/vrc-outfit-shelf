DROP TABLE IF EXISTS compatibility;
DROP TABLE IF EXISTS preset_items;
DROP TABLE IF EXISTS assets;
DROP TABLE IF EXISTS avatar_mods;
DROP TABLE IF EXISTS presets;
DROP TABLE IF EXISTS avatars;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS category;

-------------------------------------------------

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

CREATE TYPE category AS ENUM (
    'clothing', 'hair', 'accessory', 'prop', 'texture', 'tool', 'animation', 'other'
);

CREATE TABLE assets (
    asset_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    asset_name TEXT NOT NULL,
    creator TEXT NOT NULL,
    shop_url TEXT,
    memo TEXT,
    asset_category category NOT NULL
);

CREATE TABLE avatar_mods (
    mod_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    avatar_id INT NOT NULL REFERENCES avatars(avatar_id) ON DELETE CASCADE,
    mod_description TEXT NOT NULL,
    modded_date DATE DEFAULT CURRENT_DATE
);

CREATE TABLE presets (
    preset_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    avatar_id INT NOT NULL REFERENCES avatars(avatar_id) ON DELETE CASCADE,
    preset_name TEXT NOT NULL,
    memo TEXT
);

CREATE TABLE preset_items (
    preset_id INT NOT NULL REFERENCES presets(preset_id) ON DELETE CASCADE,
    asset_id INT NOT NULL REFERENCES assets(asset_id) ON DELETE CASCADE,
    memo TEXT,
    PRIMARY KEY (preset_id, asset_id)
);

CREATE TABLE compatibility (
    avatar_id INT NOT NULL REFERENCES avatars(avatar_id) ON DELETE CASCADE,
    asset_id INT NOT NULL REFERENCES assets(asset_id) ON DELETE CASCADE,
    memo TEXT,
    PRIMARY KEY (avatar_id, asset_id)
);