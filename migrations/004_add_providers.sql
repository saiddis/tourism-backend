-- Create providers table
CREATE TABLE providers (
    id                  SERIAL PRIMARY KEY,
    user_id             INT UNIQUE NOT NULL REFERENCES users(id),
    phone               VARCHAR(20) NOT NULL,
    provider_type       VARCHAR(20) NOT NULL CHECK (provider_type IN ('private', 'organization')),
    instagram_url       TEXT,
    telegram_url        TEXT,
    facebook_url        TEXT,
    years_experience    INT DEFAULT 0,
    bio                 TEXT,
    active              BOOLEAN NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create provider_applications table
CREATE TABLE provider_applications (
    id                  SERIAL PRIMARY KEY,
    user_id             INT NOT NULL REFERENCES users(id),
    phone               VARCHAR(20) NOT NULL,
    provider_type       VARCHAR(20) NOT NULL CHECK (provider_type IN ('private', 'organization')),
    instagram_url        TEXT,
    telegram_url        TEXT,
    facebook_url        TEXT,
    years_experience    INT DEFAULT 0,
    bio                 TEXT,
    status              VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'rejected')),
    admin_note          TEXT,
    created_at          TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Add provider_id to tours (nullable - existing tours have no provider)
ALTER TABLE tours ADD COLUMN provider_id INT REFERENCES providers(id);

-- Create indexes for performance
CREATE INDEX idx_tours_provider_id ON tours(provider_id);
CREATE INDEX idx_provider_applications_user_id ON provider_applications(user_id);
CREATE INDEX idx_provider_applications_status ON provider_applications(status);
CREATE INDEX idx_providers_active ON providers(active);
CREATE INDEX idx_providers_user_id ON providers(user_id);
