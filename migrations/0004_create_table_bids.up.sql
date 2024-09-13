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

CREATE TABLE IF NOT EXISTS bid (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    status bid_status_type DEFAULT 'Created',
    tender_id UUID REFERENCES tender(id) ON DELETE CASCADE,
    author_type bid_author_type,
    author_id UUID NOT NULL,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS review (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bid_id UUID REFERENCES bid(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)