INSERT INTO users (vrc_name, password_hash) VALUES
    ('vrc_name', 'example_hash');
INSERT INTO avatars (user_id, avatar_name, creator, shop_url, memo) VALUES
    (1, 'avatar_name1', 'creator_name1', 'example_url', 'hoge'),
    (1, 'avatar_name2', 'creator_name1', NULL, 'hoge');
INSERT INTO assets (user_id, asset_name, creator, shop_url, memo, asset_category) VALUES
    (1, 'example_asset1', 'creator_name2', 'example_url', 'hoge', 'clothing'),
    (1, 'example_asset2', 'creator_name3', 'example_url', 'hoge', 'hair'),
    (1, 'example_asset3', 'creator_name4', 'example_url', 'hoge', 'accessory');
INSERT INTO compatibility (avatar_id, asset_id, memo) VALUES
    (1, 1, 'hoge'),
    (1, 2, 'hoge'),
    (2, 1, 'hoge'),
    (2, 3, NULL);
INSERT INTO presets (avatar_id, preset_name, memo) VALUES
    (1, 'example_preset1', 'hoge'),
    (1, 'example_preset2', 'hoge'),
    (2, 'example_preset3', 'hoge');
INSERT INTO preset_items (preset_id, asset_id, memo) VALUES
    (1, 1, 'hoge'), (1, 2, 'hoge'),
    (2, 2, 'hoge'),
    (3, 1, 'hoge');