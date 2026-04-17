package queries

// TOURS
const (
	CreateTour = `
INSERT INTO tours (destination_id,provider_id,name,description,price,start_date,end_date,capacity)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
returning id,created_at`
	GetTourByID = `
SELECT
	t.id,
	t.destination_id,
	t.provider_id,
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
FROM tours t
JOIN destinations d ON d.id = t.destination_id
WHERE t.id = $1`
	GetAllTours = `
SELECT
	t.id,
	t.destination_id,
	t.provider_id,
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
FROM tours t
JOIN destinations d ON d.id = t.destination_id
ORDER BY t.start_date, t.id`
	UpdateTour = `
UPDATE tours 
	SET destination_id=$1, provider_id=$2, name=$3, description=$4,
		price=$5, start_date=$6, end_date=$7, capacity=$8
	WHERE id=$9
	RETURNING created_at`
	DeleteTour = `
	DELETE FROM tours WHERE id = $1`
	GetToursByDestinationID = `
SELECT
	t.id,
	t.destination_id,
	t.provider_id,
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
FROM tours t
JOIN destinations d ON d.id = t.destination_id
WHERE t.destination_id = $1
ORDER BY t.start_date, t.id`
	DecrementTourCapacity = `
UPDATE tours SET capacity = capacity - 1 WHERE id = $1 AND capacity > 0`
	IncrementTourCapacity = `
UPDATE tours SET capacity = capacity + 1 WHERE id = $1`
)

// TOUR HIGHLIGHTS
const (
	CreateTourHighlight = `
INSERT INTO tour_highlights (tour_id, image_url, title, sort_order)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at`
	GetTourHighlightsByTourID = `
SELECT id, tour_id, image_url, title, sort_order, created_at
FROM tour_highlights
WHERE tour_id = $1
ORDER BY sort_order, id`
	DeleteTourHighlight = `DELETE FROM tour_highlights WHERE id = $1`
)
