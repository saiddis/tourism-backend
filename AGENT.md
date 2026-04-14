# Frontend Agent Guide

This file is for agents working on the frontend that consumes this backend.

## Project Summary

- Backend stack: Go, `chi`, PostgreSQL.
- Default API base URL: `http://localhost:8080`.
- API style: unversioned REST-like JSON endpoints.
- Response format:
  - Success: JSON object or JSON array.
  - Error: `{ "error": "message" }`.
- Timestamps are Go `time.Time` JSON values. Treat them as ISO/RFC3339 strings in the frontend.

## Current Backend Limits

Build the frontend around the backend that exists now, not the one you wish existed.

- Login and token refresh exist, but there is still no logout or current-user endpoint.
- JWT auth is enforced on protected routes with access tokens, and role checks exist for manager/admin endpoints.
- No CORS middleware is configured.
- No pagination, sorting, filtering, or search endpoints exist.
- Delete/status-update endpoints usually return `{ "message": "..." }`, not the updated resource.
- Some service methods exist but are not exposed as routes.

## Frontend Integration Rules

- Keep one API client module for all requests.
- Keep snake_case at the network boundary because backend JSON uses snake_case.
- If the UI prefers camelCase, map once in the API layer instead of mixing naming styles across components.
- Always parse non-2xx responses as JSON first and surface `error` when present.
- Expect raw database constraint errors in some failure cases. Show a user-friendly fallback if the message is not clean.
- Because CORS is not configured, use a frontend dev proxy or serve frontend/backend from the same origin in development.
- Auth responses now return access and refresh tokens. Persist both if the frontend needs session restore.

## API Contract

### Users

- `POST /auth/register`
  - Body:
    ```json
    {
      "name": "John",
      "email": "john@example.com",
      "password": "secret"
    }
    ```
  - Returns: created user object.
  - Notes:
    - `role` is assigned by backend as `client`.
    - Password is not returned.

- `POST /auth/login`
  - Body:
    ```json
    {
      "email": "john@example.com",
      "password": "secret"
    }
    ```
  - Returns:
    ```json
    {
      "access_token": "jwt",
      "refresh_token": "jwt",
      "token_type": "Bearer",
      "access_expires_in": 900,
      "refresh_expires_in": 604800
    }
    ```

- `POST /auth/refresh`
  - Body:
    ```json
    {
      "refresh_token": "jwt"
    }
    ```
  - Returns: same token pair shape as login.

- `GET /users`
  - Returns: `User[]`

- `GET /users/{id}`
  - Returns: `User`

### Destinations

- `POST /destinations`
  - Body:
    ```json
    {
      "name": "Paris, France",
      "description": "The city of light, love, and iconic landmarks",
      "image_url": "https://images.unsplash.com/photo-1502602898657-3e91760cbb34?w=800"
    }
    ```
  - Returns: created destination

- `GET /destinations`
  - Returns: `Destination[]`

- `GET /destinations/{id}`
  - Returns: `Destination`

- `DELETE /destinations/{id}`
  - Returns:
    ```json
    { "message": "destination deleted" }
    ```

### Tours

- `POST /tours`
  - Body:
    ```json
    {
      "destination_id": 1,
      "name": "Weekend Tour",
      "description": "Short trip",
      "price": 199.99,
      "start_date": "2026-05-01T09:00:00Z",
      "end_date": "2026-05-03T18:00:00Z",
      "capacity": 20
    }
    ```
  - Returns: created tour

- `GET /tours`
  - Returns: `Tour[]`

- `GET /tours/{id}`
  - Returns: `Tour`

- `PUT /tours/{id}`
  - Body: same shape as create
  - Returns: updated tour

- `DELETE /tours/{id}`
  - Returns:
    ```json
    { "message": "tour deleted" }
    ```

### Bookings

- `POST /bookings`
  - Body:
    ```json
    {
      "user_id": 1,
      "tour_id": 2
    }
    ```
  - Returns: created booking
  - Notes:
    - Backend takes `user_id` from the authenticated access token and ignores any mismatched body value.
    - Backend forces `status` to `pending`, even if client sends another value.

- `GET /bookings`
  - Returns: `Booking[]`

- `GET /bookings/user/{id}`
  - Returns: bookings for a user
  - Notes:
    - Clients can fetch only their own bookings.
    - Managers and admins can fetch bookings for any user id.

- `PUT /bookings/{id}/status`
  - Body:
    ```json
    { "status": "confirmed" }
    ```
  - Allowed values in code: `pending`, `confirmed`, `cancelled`
  - Returns:
    ```json
    { "message": "status updated" }
    ```

- `DELETE /bookings/{id}`
  - Returns:
    ```json
    { "message": "booking deleted" }
    ```

### Payments

- `POST /payments`
  - Body:
    ```json
    {
      "booking_id": 1,
      "amount": 199.99
    }
    ```
  - Returns: created payment
  - Notes:
    - Backend forces `status` to `pending`.
    - Clients can create payments only for their own bookings.
    - Managers and admins can create payments for any booking id.

- `GET /payments/{id}`
  - Returns: `Payment`

- `PUT /payments/{id}/status`
  - Body:
    ```json
    { "status": "paid" }
    ```
  - Allowed values in code: `pending`, `paid`, `failed`
  - Returns:
    ```json
    { "message": "payment status updated" }
    ```

### Reviews

- `POST /reviews`
  - Body:
    ```json
    {
      "user_id": 1,
      "tour_id": 2,
      "rating": 5,
      "comment": "Excellent"
    }
    ```
  - Returns: created review
  - Notes:
    - Backend takes `user_id` from the authenticated access token and ignores any mismatched body value.

- `GET /reviews/tour/{id}`
  - Returns: reviews for a tour

- `DELETE /reviews/{id}`
  - Returns:
    ```json
    { "message": "review deleted" }
    ```

## Wire Models

Use these exact API field names at the HTTP boundary.

```ts
type UserRole = 'client' | 'manager' | 'admin';

type User = {
  id: number;
  name: string;
  email: string;
  role: UserRole;
  created_at: string;
};

type Destination = {
  id: number;
  name: string;
  description: string;
  image_url: string;
  created_at: string;
};

type Tour = {
  id: number;
  destination_id: number;
  name: string;
  description: string;
  price: number;
  start_date: string;
  end_date: string;
  capacity: number;
  created_at: string;
};

type BookingStatus = 'pending' | 'confirmed' | 'cancelled';

type Booking = {
  id: number;
  user_id: number;
  tour_id: number;
  status: BookingStatus;
  created_at: string;
};

type PaymentStatus = 'pending' | 'paid' | 'failed';

type Payment = {
  id: number;
  booking_id: number;
  amount: number;
  status: PaymentStatus;
  created_at: string;
};

type Review = {
  id: number;
  user_id: number;
  tour_id: number;
  rating: number;
  comment: string;
  created_at: string;
};
```

## Missing or Non-Exposed Backend Features

Do not invent frontend flows that depend on endpoints that do not exist.

- No logout endpoint.
- No current-user endpoint like `/me`.
- No `GET /bookings/{id}` route.
- No `GET /payments/by-booking/{bookingId}` route.
- No route to fetch tours by destination directly, even though repository support exists.
- No review update endpoint.
- No destination update endpoint.
- No user update endpoint exposed in router.

## Recommended Frontend Scope

Good first-pass UI flows for this backend:

- destination list
- tour list
- tour detail
- user registration
- create booking
- user bookings
- create payment
- tour reviews

## UX Guidance

- On create/update/delete, refetch the affected list or detail unless you have a very small local state update.
- For status changes, update UI from the requested status only after the request succeeds.
- Add empty states because list endpoints may legitimately return `[]`.
- Add defensive loading and error states everywhere. This backend does not wrap responses in a consistent envelope.

## Local Development Assumptions

- Backend server port defaults to `8080`.
- Database is PostgreSQL.
- Useful backend env vars:
  - `DB_HOST`
  - `DB_PORT`
  - `DB_USER`
  - `DB_PASSWORD`
  - `DB_NAME`
  - `SERVER_PORT`
  - `ACCESS_TOKEN_SECRET`
  - `REFRESH_TOKEN_SECRET`

## Source of Truth

If backend behavior and frontend assumptions disagree, trust the backend code in:

- `internal/handler`
- `internal/service`
- `internal/domain`
- `migrations/001_init.sql`

Update this file when routes, payloads, or auth behavior change.
