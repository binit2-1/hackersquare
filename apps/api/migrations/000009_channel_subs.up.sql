CREATE TABLE channel_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    platform VARCHAR(50) NOT NULL, -- 'discord' or 'telegram'
    channel_id VARCHAR(255) NOT NULL UNIQUE, -- The Discord Channel ID or Telegram Chat ID
    chat_id VARCHAR(255) NOT NULL UNIQUE, -- The Discord Channel ID or Telegram Chat ID (same as channel_id for Telegram)
    tech_tags TEXT[] DEFAULT '{}', -- e.g., {"Web3", "Solana"}
    country VARCHAR(100), -- e.g., "India"
    is_active BOOLEAN DEFAULT true
);

CREATE INDEX idx_channel_subscriptions_chat_id ON channel_subscriptions(chat_id);