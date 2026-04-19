-- ============================================================================
-- Default Admin User (password: admin123)
-- ============================================================================
INSERT INTO users (name, email, password, role, avatar_url)
VALUES ('said', 'saiddisruptive@gmail.com', '$2a$10$LEk21wzthcBl6FjBHB/pVOKxn8vRSToZjYE4MZwzuXIhJYjZ1IUvK', 'admin', '/uploads/avatars/avatar_1.png')
ON CONFLICT (email) DO NOTHING;

-- ============================================================================
-- Provider User (password: shon123)
-- ============================================================================
INSERT INTO users (name, email, password, role)
VALUES ('shon', 'raimdodov.sh@gmail.com', '$2a$10$HyKedJv4ztCJ70dRKzj5x.mrbAfZizMHyl.mCLY4SgZvwVw2ojKGy', 'provider')
ON CONFLICT (email) DO NOTHING;

-- ============================================================================
-- Destinations
-- ============================================================================
INSERT INTO destinations (name, description, image_url, created_at)
VALUES 
    ('Dushanbe', 'Monday', '/uploads/destinations/destination_1776603928761717545.jpg', now()),
    ('Safed Dara', 'Snow, mountains and snow', '/uploads/destinations/destination_1776532982171509119.jpg', now()),
    ('Varzob', 'Mountains, rivers, summer', '/uploads/destinations/destination_1776533323259915456.jpg', now()),
    ('Pamir', 'Shirchay', '/uploads/destinations/destination_1776533770268946316.jpg', now()),
    ('Panjakent', 'Lakes and more', '/uploads/destinations/destination_1776599441875234835.jpg', now()),
    ('Sughd Province', 'Leninabad', '/uploads/destinations/destination_1776616606756502791.jpg', now())

ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Providers
-- ============================================================================
INSERT INTO providers (user_id, phone, provider_type, years_experience, bio, active, created_at)
VALUES 
    (2, '777777777777', 'private', 1, 'no', true, now())
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Tours
-- ============================================================================
INSERT INTO tours (destination_id, provider_id, name, description, price, start_date, end_date, capacity, created_at)
VALUES 
    (4, 1, 'Wakhan Walley', 'Somewhere in Pamir', 149.00, now() + interval '2 days', now() + interval '3 days', 4, now()),
    (5, 1, 'Seven Lakes', 'Panjakent to Seven Lakes for day trip and swim', 400.00, now() + interval '3 days', now() + interval '6 days', 5, now()),
    (6, 1, 'Iskandarkul Lake', 'Never been there? Me too.', 899.00, now() + interval '3 days', now() + interval '4 days', 2, now())
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Tour Highlights
-- ============================================================================
INSERT INTO tour_highlights (tour_id, image_url, created_at) VALUES
    (1, '/uploads/highlights/highlight_1_2.jpg', now()),
    (2, '/uploads/highlights/highlight_2_2.jpg', now()),
    (3, '/uploads/highlights/highlight_3_2.jpg', now())
ON CONFLICT (id) DO NOTHING;
