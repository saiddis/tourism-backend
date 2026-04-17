-- Create tour_highlights table
CREATE TABLE IF NOT EXISTS tour_highlights (
    id          SERIAL PRIMARY KEY,
    tour_id     INT NOT NULL REFERENCES tours(id) ON DELETE CASCADE,
    image_url   TEXT NOT NULL,
    title       TEXT,
    sort_order  INT DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_tour_highlights_tour_id ON tour_highlights(tour_id);
CREATE INDEX idx_tour_highlights_sort_order ON tour_highlights(tour_id, sort_order);