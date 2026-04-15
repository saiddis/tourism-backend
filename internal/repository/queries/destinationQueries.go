package queries

const (
	CreateDestination = `
INSERT INTO destinations (name, description, image_url)
VALUES ($1, $2, $3)
RETURNING id, image_url, created_at`
	GetDestinationByID = `
SELECT id, name, description, image_url, created_at
FROM destinations WHERE id = $1`
	GetAllDestinations = `
SELECT id, name, description, image_url, created_at
FROM destinations
ORDER BY name`
	DeleteDestination = `
DELETE FROM destinations WHERE id = $1`
)
