CREATE TABLE sentiments (
    id VARCHAR PRIMARY KEY,
    document_id INT NOT NULL,
    excerpt VARCHAR NOT NULL,
    comment TEXT NOT NULL,
    emotion VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL
);