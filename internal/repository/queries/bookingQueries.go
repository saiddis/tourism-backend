package queries

const (
	CreateBooking = `
INSERT INTO bookings (user_id, tour_id, status)
VALUES ($1, $2, $3)
returning id,created_at`
	GetBookingByID = `
SELECT b.id,b.user_id,b.tour_id,b.status,b.created_at,
	t.id,t.destination_id,t.name,t.description,t.price,t.start_date,t.end_date,t.capacity,t.created_at,
	d.id,d.name,d.description,d.image_url,d.created_at,
	COALESCE(
		(SELECT t.capacity - COALESCE(SUM(CASE WHEN b2.status IN ('pending', 'confirmed') THEN 1 ELSE 0 END), 0)
		FROM bookings b2 WHERE b2.tour_id = t.id GROUP BY t.capacity),
		t.capacity
	) as remaining_spots
FROM bookings b
JOIN tours t ON t.id = b.tour_id
JOIN destinations d ON d.id = t.destination_id
WHERE b.id = $1`
	GetBookingsByUserID = `
SELECT b.id,b.user_id,b.tour_id,b.status,b.created_at,
	t.id,t.destination_id,t.name,t.description,t.price,t.start_date,t.end_date,t.capacity,t.created_at,
	d.id,d.name,d.description,d.image_url,d.created_at,
	COALESCE(
		(SELECT t.capacity - COALESCE(SUM(CASE WHEN b2.status IN ('pending', 'confirmed') THEN 1 ELSE 0 END), 0)
		FROM bookings b2 WHERE b2.tour_id = t.id GROUP BY t.capacity),
		t.capacity
	) as remaining_spots
FROM bookings b
JOIN tours t ON t.id = b.tour_id
JOIN destinations d ON d.id = t.destination_id
WHERE b.user_id = $1
ORDER BY b.created_at DESC`
	GetAllBookings = `
SELECT b.id,b.user_id,b.tour_id,b.status,b.created_at,
	t.id,t.destination_id,t.name,t.description,t.price,t.start_date,t.end_date,t.capacity,t.created_at,
	d.id,d.name,d.description,d.image_url,d.created_at,
	COALESCE(
		(SELECT t.capacity - COALESCE(SUM(CASE WHEN b2.status IN ('pending', 'confirmed') THEN 1 ELSE 0 END), 0)
		FROM bookings b2 WHERE b2.tour_id = t.id GROUP BY t.capacity),
		t.capacity
	) as remaining_spots
FROM bookings b
JOIN tours t ON t.id = b.tour_id
JOIN destinations d ON d.id = t.destination_id
ORDER BY b.created_at DESC`
	UpdateBookingStatus             = `UPDATE bookings SET status = $1 WHERE id = $2`
	DeleteBooking                   = `DELETE FROM bookings WHERE id = $1`
	GetPendingBookingsCountByTourID = `
SELECT COUNT(*) FROM bookings WHERE tour_id = $1 AND status = 'pending'`
	GetActiveBookingsCountByTourID = `
SELECT COUNT(*) FROM bookings WHERE tour_id = $1 AND status IN ('pending', 'confirmed')`
	GetPendingBookingsByTourID = `
SELECT b.id, b.user_id, b.tour_id, b.status, b.created_at, u.email, u.name
FROM bookings b
JOIN users u ON u.id = b.user_id
WHERE b.tour_id = $1 AND b.status = 'pending'
ORDER BY b.created_at`
	ConfirmPendingBookingsForTour = `
UPDATE bookings SET status = 'confirmed' WHERE tour_id = $1 AND status = 'pending'`
	GetToursNeedingConfirmation = `
SELECT t.id, t.name, t.start_date, t.capacity, 
	(SELECT COUNT(*) FROM bookings WHERE tour_id = t.id AND status = 'pending') as pending_count
FROM tours t
WHERE t.start_date <= NOW() + INTERVAL '24 hours'
	AND EXISTS (SELECT 1 FROM bookings WHERE tour_id = t.id AND status = 'pending')
ORDER BY t.start_date`
	MarkBookingsCompleted = `
UPDATE bookings 
SET status = 'completed' 
WHERE status = 'confirmed' 
AND tour_id IN (SELECT id FROM tours WHERE end_date < NOW())`
	GetCompletedBookingsByUserID = `
SELECT b.id,b.user_id,b.tour_id,b.status,b.created_at,
	t.id,t.destination_id,t.name,t.description,t.price,t.start_date,t.end_date,t.capacity,t.created_at,
	d.id,d.name,d.description,d.image_url,d.created_at,
	COALESCE(
		(SELECT t.capacity - COALESCE(SUM(CASE WHEN b2.status IN ('pending', 'confirmed') THEN 1 ELSE 0 END), 0)
		FROM bookings b2 WHERE b2.tour_id = t.id GROUP BY t.capacity),
		t.capacity
	) as remaining_spots
FROM bookings b
JOIN tours t ON t.id = b.tour_id
JOIN destinations d ON d.id = t.destination_id
WHERE b.user_id = $1 AND b.status = 'completed'
ORDER BY b.created_at DESC`
)
