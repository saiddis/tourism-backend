package queries

const (
	CreateProviderApplication = `
INSERT INTO provider_applications (user_id, phone, provider_type, instagram_url, telegram_url, facebook_url, years_experience, bio, status, admin_token)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, created_at, updated_at`

	GetProviderApplicationByID = `
SELECT pa.id, pa.user_id, pa.phone, pa.provider_type, pa.instagram_url, pa.telegram_url, pa.facebook_url, pa.years_experience, pa.bio, pa.status, pa.admin_note, pa.admin_token, pa.created_at, pa.updated_at,
       u.name as user_name, u.email as user_email
FROM provider_applications pa
JOIN users u ON pa.user_id = u.id
WHERE pa.id = $1`

	GetProviderApplicationByUserID = `
SELECT pa.id, pa.user_id, pa.phone, pa.provider_type, pa.instagram_url, pa.telegram_url, pa.facebook_url, pa.years_experience, pa.bio, pa.status, pa.admin_note, pa.admin_token, pa.created_at, pa.updated_at,
       u.name as user_name, u.email as user_email
FROM provider_applications pa
JOIN users u ON pa.user_id = u.id
WHERE pa.user_id = $1
ORDER BY pa.created_at DESC`

	GetAllProviderApplications = `
SELECT pa.id, pa.user_id, pa.phone, pa.provider_type, pa.instagram_url, pa.telegram_url, pa.facebook_url, pa.years_experience, pa.bio, pa.status, pa.admin_note, pa.admin_token, pa.created_at, pa.updated_at,
       u.name as user_name, u.email as user_email
FROM provider_applications pa
JOIN users u ON pa.user_id = u.id
ORDER BY pa.created_at DESC`

	GetPendingProviderApplications = `
SELECT pa.id, pa.user_id, pa.phone, pa.provider_type, pa.instagram_url, pa.telegram_url, pa.facebook_url, pa.years_experience, pa.bio, pa.status, pa.admin_note, pa.admin_token, pa.created_at, pa.updated_at,
       u.name as user_name, u.email as user_email
FROM provider_applications pa
JOIN users u ON pa.user_id = u.id
WHERE pa.status = 'pending'
ORDER BY pa.created_at ASC`

	GetProviderApplicationByToken = `
SELECT pa.id, pa.user_id, pa.phone, pa.provider_type, pa.instagram_url, pa.telegram_url, pa.facebook_url, pa.years_experience, pa.bio, pa.status, pa.admin_note, pa.admin_token, pa.created_at, pa.updated_at,
       u.name as user_name, u.email as user_email
FROM provider_applications pa
JOIN users u ON pa.user_id = u.id
WHERE pa.admin_token = $1`

	UpdateProviderApplicationStatus = `
UPDATE provider_applications SET
    status = $2,
    admin_note = $3,
    updated_at = NOW()
WHERE id = $1`
)
