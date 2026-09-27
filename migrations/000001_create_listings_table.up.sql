CREATE TABLE listings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    price BIGINT NOT NULL,
    image_url TEXT NOT NULL,
    city TEXT NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Indexes
-- CREATE INDEX idx_listings_title ON listings(title);
-- CREATE INDEX idx_listings_price ON listings(price);
-- CREATE INDEX idx_listings_created_at ON listings(created_at);
