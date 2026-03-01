CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE hackathons (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    host TEXT NOT NULL,
    location TEXT,
    prize_usd NUMERIC(10, 2),
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,

    --FTS 
    search_vector tsvector GENERATED ALWAYS AS (
        to_tsvector('english', coalesce(title, '') || ' ' || coalesce(host, '') || ' ' || coalesce(location, ''))
    ) STORED,


    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

--create GIN index on the search vector 
CREATE INDEX hackathons_search_idx ON hackathons USING GIN (search_vector);