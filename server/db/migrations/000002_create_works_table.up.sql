CREATE TABLE IF NOT EXISTS works (
    id INTEGER,
    user_id INTEGER,
    title VARCHAR(256),
    work_url VARCHAR(256),
    count INTEGER,
    bookmark BOOLEAN,
    created_at DATE,
    updated_at DATE
);