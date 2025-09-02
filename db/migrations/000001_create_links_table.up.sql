CREATE TABLE links (
    short VARCHAR(20) PRIMARY KEY,
    original TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
