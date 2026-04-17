package queries

// TOURS
const (
	CreateTour = `
INSERT INTO tours (destination_id,name,description,price,start_date,end_date,capacity)
VALUES ($1,$2,$3,$4,$5,$6,$7)
returning id,created_at`
	GetTourByID = `
SELECT
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
FROM tours t
JOIN destinations d ON d.id = t.destination_id
WHERE t.id = $1`
	GetAllTours = `
SELECT
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
FROM tours t
JOIN destinations d ON d.id = t.destination_id
ORDER BY t.start_date, t.id`
	UpdateTour = `
UPDATE tours 
	SET destination_id=$1, name=$2, description=$3,
		price=$4, start_date=$5, end_date=$6, capacity=$7
	WHERE id=$8
	RETURNING created_at`
	DeleteTour = `
	DELETE FROM tours WHERE id = $1`
	GetToursByDestinationID = `
SELECT
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
FROM tours t
JOIN destinations d ON d.id = t.destination_id
WHERE t.destination_id = $1
ORDER BY t.start_date, t.id`
	DecrementTourCapacity = `
UPDATE tours SET capacity = capacity - 1 WHERE id = $1 AND capacity > 0`
	IncrementTourCapacity = `
UPDATE tours SET capacity = capacity + 1 WHERE id = $1`
)
