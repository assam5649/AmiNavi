CREATE TABLE IF NOT EXISTS users (
    id INTEGER,
    firebase_uid VARCHAR(256),
    login_id VARCHAR(256),
    display_name VARCHAR(256),
    profile_image_url VARCHAR(256),
    email VARCHAR(256),
    created_at DATE,
    updated_at DATE
)