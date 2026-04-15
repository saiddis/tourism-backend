package queries

// USER
const (
	CreateUser = `
INSERT INTO users (name,email,password,role)
VALUES ($1,$2,$3,$4)
RETURNING id,created_at`
	GetUserByEmail = `
SELECT id, name, email, password, role, created_at 
FROM users WHERE email = $1`
	GetAllUsers = `
SELECT id, name, email, role, created_at FROM users`
	UpdateUser = `
UPDATE users SET name = $1, email = $2 where id = $3`
	UpdateUserPassword = `
UPDATE users SET password = $1 WHERE id = $2`
	DeleteUser = `
DELETE FROM users WHERE id = $1`
	GetUserByID = `
SELECT id, name, email, role, created_at
FROM users WHERE id = $1`
)
