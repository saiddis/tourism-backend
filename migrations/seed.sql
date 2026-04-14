INSERT INTO users (name, email, password, role) VALUES
('John Doe', 'john@example.com', 'password123', 'client'),
('Jane Smith', 'jane@example.com', 'password123', 'client'),
('Admin User', 'admin@tourism.tj', 'admin123', 'admin');

INSERT INTO destinations (name, description, image_url) VALUES
('Paris, France', 'The city of light, love, and iconic landmarks', 'https://images.unsplash.com/photo-1502602898657-3e91760cbb34?w=800'),
('Tokyo, Japan', 'A vibrant metropolis where tradition meets innovation', 'https://images.unsplash.com/photo-1540959733332-eab4deabeeaf?w=800'),
('Santorini, Greece', 'Stunning sunsets and white-washed buildings overlooking the Aegean Sea', 'https://images.unsplash.com/photo-1613395877344-13d4a8e0d49e?w=800');

INSERT INTO tours (destination_id, name, description, price, start_date, end_date, capacity) VALUES
(1, 'Paris City Lights', 'Experience the magic of Paris from the Eiffel Tower to charming cafes', 649.99, '2026-05-15 09:00:00', '2026-05-20 18:00:00', 20),
(2, 'Tokyo Discovery', 'Explore ancient temples, modern districts, and Japanese cuisine', 2199.99, '2026-06-10 08:00:00', '2026-06-17 20:00:00', 15),
(3, 'Santorini Sunset', 'Relax on beautiful beaches and watch breathtaking sunsets', 949.99, '2026-03-01 10:00:00', '2026-03-07 16:00:00', 18);

INSERT INTO bookings (user_id, tour_id, status) VALUES
(1, 1, 'confirmed' ),
(1, 2, 'pending'),
(2, 3, 'cancelled');

-- INSERT INTO payments (booking_id, amount, status) VALUES
-- (1, 1299.99, 'paid'),
-- (2, 2199.99, 'pending'),
-- (3, 1899.99, 'refunded');
