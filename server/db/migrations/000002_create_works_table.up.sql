CREATE INDEX idx_users_firebase_uid ON users (firebase_uid);

CREATE TABLE IF NOT EXISTS works (
    id INT AUTO_INCREMENT PRIMARY KEY NOT NULL,
    author VARCHAR(256)
        CHARACTER SET utf8mb4
        COLLATE utf8mb4_0900_ai_ci
        NOT NULL,
    title VARCHAR(256),
    file_name VARCHAR(256),
    raw_index INT,
    stitch_index INT,
    bookmark BOOLEAN,
    is_completed BOOLEAN,
    description VARCHAR(256),
    completed_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_works_author (author),

    FOREIGN KEY (author) REFERENCES users(firebase_uid)
        ON UPDATE CASCADE
        ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
