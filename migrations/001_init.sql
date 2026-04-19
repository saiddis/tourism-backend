-- ============================================================================
-- Tourism Platform Database Schema
-- Combined migration file for PostgreSQL
-- ============================================================================

-- ============================================================================
-- Users Table
-- ============================================================================
CREATE TABLE IF NOT EXISTS users (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    email      VARCHAR(100) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    role       VARCHAR(20)  NOT NULL DEFAULT 'client',
    avatar_url VARCHAR(500),
    balance   DECIMAL(12,2) DEFAULT 0.00,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- Destinations Table
-- ============================================================================
CREATE TABLE IF NOT EXISTS destinations (
    id                SERIAL PRIMARY KEY,
    name              VARCHAR(100) NOT NULL UNIQUE,
    description       TEXT,
    image_url         TEXT,
    created_at        TIMESTAMP NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- Providers Table
-- ============================================================================
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


-- ============================================================================
-- Tours Table
-- ============================================================================
CREATE TABLE IF NOT EXISTS tours (
    id             SERIAL PRIMARY KEY,
    destination_id INT            NOT NULL REFERENCES destinations(id),
    provider_id    INT            REFERENCES providers(id),
    name           VARCHAR(200)   NOT NULL,
    description    TEXT,
    price          DECIMAL(10,2)  NOT NULL,
    start_date     TIMESTAMP      NOT NULL,
    end_date       TIMESTAMP      NOT NULL,
    capacity       INT            NOT NULL,
    created_at     TIMESTAMP      NOT NULL DEFAULT NOW()
);

-- Add provider_id column if providers table doesn't exist yet (handled by 004 migration)
-- This will be applied after providers table is created

-- ============================================================================
-- Bookings Table
-- ============================================================================
CREATE TABLE IF NOT EXISTS bookings (
    id         SERIAL PRIMARY KEY,
    user_id    INT          NOT NULL REFERENCES users(id),
    tour_id    INT          NOT NULL REFERENCES tours(id),
    status     VARCHAR(20)  NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- Payments Table
-- ============================================================================
CREATE TABLE IF NOT EXISTS payments (
    id         SERIAL PRIMARY KEY,
    booking_id INT           NOT NULL REFERENCES bookings(id),
    amount     DECIMAL(10,2) NOT NULL,
    currency   VARCHAR(3)   DEFAULT 'TJS',
    status     VARCHAR(20)   NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP     NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- Reviews Table
-- ============================================================================
CREATE TABLE IF NOT EXISTS reviews (
    id         SERIAL PRIMARY KEY,
    user_id    INT       NOT NULL REFERENCES users(id),
    tour_id    INT       NOT NULL REFERENCES tours(id),
    rating     INT       NOT NULL,
    comment    TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- Provider Applications Table
-- ============================================================================
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
    admin_token         VARCHAR(64) UNIQUE,
    admin_note          TEXT,
    created_at          TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- Tour Highlights Table
-- ============================================================================
CREATE TABLE IF NOT EXISTS tour_highlights (
    id          SERIAL PRIMARY KEY,
    tour_id     INT NOT NULL REFERENCES tours(id) ON DELETE CASCADE,
    image_url   TEXT NOT NULL,
    title       TEXT,
    sort_order  INT DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- Indexes for Performance
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_tours_provider_id ON tours(provider_id);
CREATE INDEX IF NOT EXISTS idx_provider_applications_user_id ON provider_applications(user_id);
CREATE INDEX IF NOT EXISTS idx_provider_applications_status ON provider_applications(status);
CREATE INDEX IF NOT EXISTS idx_provider_applications_admin_token ON provider_applications(admin_token);
CREATE INDEX IF NOT EXISTS idx_providers_active ON providers(active);
CREATE INDEX IF NOT EXISTS idx_providers_user_id ON providers(user_id);
CREATE INDEX IF NOT EXISTS idx_tour_highlights_tour_id ON tour_highlights(tour_id);
CREATE INDEX IF NOT EXISTS idx_tour_highlights_sort_order ON tour_highlights(tour_id, sort_order);
