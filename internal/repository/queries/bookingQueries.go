package queries

const (
	CreateBooking = `
INSERT INTO bookings (user_id, tour_id, status)
VALUES ($1, $2, $3)
RETURNING id, created_at`

	GetBookingByID = `
SELECT
	b.id,
	b.user_id,
	b.tour_id,
	b.status,
	b.created_at,
	t.id,
	t.destination_id,
	t.name,
	t.description,
	t.price,
	t.start_date,
	t.end_date,
	t.capacity,
	t.created_at,
	d.id,
	d.name,
	d.description,
	d.image_url,
	d.created_at
FROM bookings b
JOIN tours t ON t.id = b.tour_id
JOIN destinations d ON d.id = t.destination_id
WHERE b.id = $1`

	GetBookingsByUserID = `
SELECT
	b.id,
	b.user_id,
	b.tour_id,
	b.status,
	b.created_at,
	t.id,
	t.destination_id,
	t.name,
	t.description,
	t.price,
	t.start_date,
	t.end_date,
	t.capacity,
	t.created_at,
	d.id,
	d.name,
	d.description,
	d.image_url,
	d.created_at
FROM bookings b
JOIN tours t ON t.id = b.tour_id
JOIN destinations d ON d.id = t.destination_id
WHERE b.user_id = $1
ORDER BY b.created_at DESC, b.id DESC`

	GetAllBookings = `
SELECT
	b.id,
	b.user_id,
	b.tour_id,
	b.status,
	b.created_at,
	t.id,
	t.destination_id,
	t.name,
	t.description,
	t.price,
	t.start_date,
	t.end_date,
	t.capacity,
	t.created_at,
	d.id,
	d.name,
	d.description,
	d.image_url,
	d.created_at
FROM bookings b
JOIN tours t ON t.id = b.tour_id
JOIN destinations d ON d.id = t.destination_id
ORDER BY b.created_at DESC, b.id DESC`

	UpdateBookingStatus = `
UPDATE bookings SET status = $1 WHERE id = $2`

	DeleteBooking = `
DELETE FROM bookings WHERE id = $1`
)
