# Tourism Backend API

Go REST API for the tourism project.

This backend uses:

- Go + Chi
- PostgreSQL
- JWT access tokens
- HttpOnly refresh token cookies

## Overview

Base URL in local development:

```text
http://localhost:8080
```

Core features:

- user registration and login
- JWT authorization with access and refresh tokens
- role-based access control for `client`, `manager`, and `admin`
- destinations, tours, bookings, payments, and reviews APIs
- flattened tour and booking responses with destination metadata included

## Authentication

The API uses a two-token flow:

- access token: returned in JSON, send it as `Authorization: Bearer <token>`
- refresh token: stored as an `HttpOnly` cookie named `refresh_token`

Important details:

- `POST /auth/login` returns the access token in JSON and sets the refresh token cookie
- `POST /auth/refresh` reads the refresh token from the cookie only
- `DELETE /auth/logout` clears the refresh token cookie
- browser clients must send cookies with `credentials: 'include'`

Refresh cookie settings in the current implementation:

- `HttpOnly`
- `Secure`
- `SameSite=None`
- `Path=/`

Access token TTL:

- 15 minutes

Refresh token TTL:

- 7 days

## Roles

- `client`: regular authenticated user
- `manager`: content and operations management
- `admin`: full access, including user administration

## Error Format

Errors are returned as JSON:

```json
{
  "error": "message"
}
```

## Environment Variables

The server reads:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=tourism_db
SERVER_PORT=8080
ACCESS_TOKEN_SECRET=change-me
REFRESH_TOKEN_SECRET=change-me-too
```

## Running Locally

1. Start PostgreSQL and apply migrations.
2. Set the environment variables.
3. Run the server:

```bash
go run ./cmd
```

With Docker Compose:

```bash
docker compose up --build
```

## Data Models

### User

```json
{
  "id": 1,
  "name": "Jane Doe",
  "email": "jane@example.com",
  "role": "client",
  "created_at": "2026-04-15T10:00:00Z"
}
```

### Destination

```json
{
  "id": 1,
  "name": "Dushanbe",
  "description": "Capital city of Tajikistan",
  "image_url": "https://images.unsplash.com/...",
  "created_at": "2026-04-15T10:00:00Z"
}
```

### Tour

`GET` tour responses include flattened destination fields:

```json
{
  "id": 1,
  "destination_id": 1,
  "destination_name": "Dushanbe",
  "destination_description": "Capital city of Tajikistan",
  "destination_image_url": "https://images.unsplash.com/...",
  "name": "Dushanbe City Highlights",
  "description": "A guided city experience",
  "price": 99.5,
  "start_date": "2026-05-01T09:00:00Z",
  "end_date": "2026-05-03T18:00:00Z",
  "capacity": 12,
  "created_at": "2026-04-15T10:00:00Z"
}
```

### Booking

`GET` booking responses include flattened tour and destination fields:

```json
{
  "id": 10,
  "user_id": 2,
  "tour_id": 1,
  "tour_name": "Dushanbe City Highlights",
  "tour_description": "A guided city experience",
  "tour_price": 99.5,
  "tour_start_date": "2026-05-01T09:00:00Z",
  "tour_end_date": "2026-05-03T18:00:00Z",
  "tour_capacity": 12,
  "destination_id": 1,
  "destination_name": "Dushanbe",
  "destination_description": "Capital city of Tajikistan",
  "destination_image_url": "https://images.unsplash.com/...",
  "status": "pending",
  "created_at": "2026-04-15T10:00:00Z"
}
```

### Payment

```json
{
  "id": 1,
  "booking_id": 10,
  "amount": 99.5,
  "status": "pending",
  "created_at": "2026-04-15T10:00:00Z"
}
```

### Review

```json
{
  "id": 1,
  "user_id": 2,
  "tour_id": 1,
  "rating": 5,
  "comment": "Excellent trip",
  "created_at": "2026-04-15T10:00:00Z"
}
```

## Endpoint Reference

### Auth

#### `POST /auth/register`

Public.

Request:

```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "password": "secret123"
}
```

Response `201 Created`:

```json
{
  "id": 2,
  "name": "Jane Doe",
  "email": "jane@example.com",
  "role": "client",
  "created_at": "2026-04-15T10:00:00Z"
}
```

#### `POST /auth/login`

Public.

Request:

```json
{
  "email": "jane@example.com",
  "password": "secret123"
}
```

Response `200 OK`:

```json
{
  "access_token": "jwt-access-token",
  "token_type": "Bearer",
  "access_expires_in": 900,
  "refresh_expires_in": 604800
}
```

Also sets:

```text
Set-Cookie: refresh_token=...; HttpOnly; Secure; SameSite=None; Path=/
```

#### `POST /auth/refresh`

Public, but requires the `refresh_token` cookie.

Request body:

```json
{}
```

Response `200 OK`:

```json
{
  "access_token": "new-jwt-access-token",
  "token_type": "Bearer",
  "access_expires_in": 900,
  "refresh_expires_in": 604800
}
```

Also rotates the refresh cookie.

#### `DELETE /auth/logout`

Public.

Response `200 OK`:

```json
{
  "message": "logged out"
}
```

### Users

#### `GET /users/me`

Authenticated.

Response `200 OK`:

```json
{
  "id": 2,
  "name": "Jane Doe",
  "email": "jane@example.com",
  "role": "client",
  "created_at": "2026-04-15T10:00:00Z"
}
```

#### `GET /users`

Admin only.

Response:

```json
[
  {
    "id": 1,
    "name": "Admin",
    "email": "admin@tourism.tj",
    "role": "admin",
    "created_at": "2026-04-15T10:00:00Z"
  }
]
```

#### `GET /users/{id}`

Admin only.

### Destinations

#### `GET /destinations`

Public.

Response:

```json
[
  {
    "id": 1,
    "name": "Dushanbe",
    "description": "Capital city of Tajikistan",
    "image_url": "https://images.unsplash.com/...",
    "created_at": "2026-04-15T10:00:00Z"
  }
]
```

#### `GET /destinations/{id}`

Public.

#### `POST /destinations`

Manager or admin.

Request:

```json
{
  "name": "Dushanbe",
  "description": "Capital city of Tajikistan",
  "image_url": "https://images.unsplash.com/..."
}
```

#### `DELETE /destinations/{id}`

Manager or admin.

Response:

```json
{
  "message": "destination deleted"
}
```

### Tours

#### `GET /tours`

Public.

Returns tours with flattened destination fields.

#### `GET /tours/{id}`

Public.

Returns a single tour with flattened destination fields.

#### `GET /tours/destination/{id}`

Public.

Returns tours for a destination.

#### `POST /tours`

Manager or admin.

Request:

```json
{
  "destination_id": 1,
  "name": "Dushanbe City Highlights",
  "description": "A guided city experience",
  "price": 99.5,
  "start_date": "2026-05-01T09:00:00Z",
  "end_date": "2026-05-03T18:00:00Z",
  "capacity": 12
}
```

#### `PUT /tours/{id}`

Manager or admin.

Request body is the same as `POST /tours`.

#### `DELETE /tours/{id}`

Manager or admin.

Response:

```json
{
  "message": "tour deleted"
}
```

### Bookings

#### `POST /bookings`

Authenticated.

The backend ignores any incoming `user_id` and uses the authenticated user.

Request:

```json
{
  "tour_id": 1
}
```

Response `201 Created`:

```json
{
  "id": 10,
  "user_id": 2,
  "tour_id": 1,
  "status": "pending",
  "created_at": "2026-04-15T10:00:00Z"
}
```

#### `GET /bookings/user/{id}`

Authenticated.

Rules:

- clients can only read their own bookings
- managers and admins can read any user's bookings

Response includes flattened `tour_*` and `destination_*` fields.

#### `PUT /bookings/{id}/status`

Authenticated.

Request:

```json
{
  "status": "cancelled"
}
```

Rules:

- clients can only cancel their own bookings
- managers and admins can set any valid booking status

Response:

```json
{
  "message": "status updated"
}
```

#### `GET /bookings`

Manager or admin.

Returns all bookings with flattened `tour_*` and `destination_*` fields.

#### `DELETE /bookings/{id}`

Manager or admin.

Response:

```json
{
  "message": "booking deleted"
}
```

### Payments

#### `POST /payments`

Authenticated.

Rules:

- clients can only create payments for their own bookings
- managers and admins can create payments for any booking

Request:

```json
{
  "booking_id": 10,
  "amount": 99.5
}
```

Response `201 Created`:

```json
{
  "id": 1,
  "booking_id": 10,
  "amount": 99.5,
  "status": "pending",
  "created_at": "2026-04-15T10:00:00Z"
}
```

#### `GET /payments/{id}`

Manager or admin.

#### `PUT /payments/{id}/status`

Manager or admin.

Request:

```json
{
  "status": "paid"
}
```

Response:

```json
{
  "message": "payment status updated"
}
```

### Reviews

#### `GET /reviews/tour/{id}`

Public.

Response:

```json
[
  {
    "id": 1,
    "user_id": 2,
    "tour_id": 1,
    "rating": 5,
    "comment": "Excellent trip",
    "created_at": "2026-04-15T10:00:00Z"
  }
]
```

#### `POST /reviews`

Authenticated.

The backend ignores any incoming `user_id` and uses the authenticated user.

Request:

```json
{
  "tour_id": 1,
  "rating": 5,
  "comment": "Excellent trip"
}
```

#### `DELETE /reviews/{id}`

Manager or admin.

Response:

```json
{
  "message": "review deleted"
}
```

## Example Authenticated Request

```bash
curl http://localhost:8080/users/me \
  -H "Authorization: Bearer <access-token>"
```

## Example Browser Refresh Flow

Frontend clients should:

1. store only the access token in app state or storage
2. send `credentials: 'include'` with API requests
3. call `POST /auth/refresh` when access token refresh is needed
4. let the browser send the `refresh_token` cookie automatically

## Notes

- list endpoints return JSON arrays
- there is no pagination yet
- create endpoints generally return the created row as stored by the write query
- flattened destination metadata is guaranteed on read endpoints for tours and bookings
