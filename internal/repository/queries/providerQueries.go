package queries

const (
	CreateProvider = `
INSERT INTO providers (user_id, phone, provider_type, instagram_url, telegram_url, facebook_url, years_experience, bio, active)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, created_at`

	GetProviderByID = `
SELECT p.id, p.user_id, p.phone, p.provider_type, p.instagram_url, p.telegram_url, p.facebook_url, p.years_experience, p.bio, p.active, p.created_at,
       u.name as user_name, u.email as user_email, u.avatar_url as user_avatar_url
FROM providers p
JOIN users u ON p.user_id = u.id
WHERE p.id = $1`

	GetProviderByUserID = `
SELECT p.id, p.user_id, p.phone, p.provider_type, p.instagram_url, p.telegram_url, p.facebook_url, p.years_experience, p.bio, p.active, p.created_at,
       u.name as user_name, u.email as user_email, u.avatar_url as user_avatar_url
FROM providers p
JOIN users u ON p.user_id = u.id
WHERE p.user_id = $1`

	GetAllProviders = `
SELECT p.id, p.user_id, p.phone, p.provider_type, p.instagram_url, p.telegram_url, p.facebook_url, p.years_experience, p.bio, p.active, p.created_at,
       u.name as user_name, u.email as user_email, u.avatar_url as user_avatar_url
FROM providers p
JOIN users u ON p.user_id = u.id
ORDER BY p.created_at DESC`

	GetActiveProviders = `
SELECT p.id, p.user_id, p.phone, p.provider_type, p.instagram_url, p.telegram_url, p.facebook_url, p.years_experience, p.bio, p.active, p.created_at,
       u.name as user_name, u.email as user_email, u.avatar_url as user_avatar_url
FROM providers p
JOIN users u ON p.user_id = u.id
WHERE p.active = true
ORDER BY p.created_at DESC`

	UpdateProvider = `
UPDATE providers SET
    phone = $2,
    provider_type = $3,
    instagram_url = $4,
    telegram_url = $5,
    facebook_url = $6,
    years_experience = $7,
    bio = $8
WHERE id = $1`

	UpdateProviderActive = `
UPDATE providers SET active = $2 WHERE id = $1`
)
