-- CREATE TYPE bid_author_type AS ENUM (
--     'Organization',
--     'User'
-- );

-- CREATE TYPE bid_status_type AS ENUM (
--     'Created',
--     'Published',
--     'Canceled',
--     'Approved',
--     'Rejected'
-- );

CREATE TABLE IF NOT EXISTS bids (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    status bid_status_type,
    tender_id INT REFERENCES tenders(id) ON DELETE CASCADE,
    author_type bid_author_type,
    author_id VARCHAR(100) NOT NULL,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);