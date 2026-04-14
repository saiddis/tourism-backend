# Tourism Platform UI Design Brief

This document is for a UI/UX designer working on the frontend for this project.

## What This Product Is

This is a tourism booking platform. Users can browse destinations and tours, register an account, create bookings, make payments, and leave reviews on tours.

The current backend is functional but minimal. The design should match what the product can actually do right now.

## Primary Product Areas

- discover destinations
- browse tours
- view tour details
- register a user account
- book a tour
- view a user's bookings
- create a payment for a booking
- read and write tour reviews

## Target Users

- Travelers looking for tours
- Managers or admins may exist in data, but there is no implemented admin/auth flow yet

For now, the most important experience is the traveler journey.

## Main Objects in the System

### Destination

A destination contains:

- name
- description
- image_url

Examples:

- Paris, France
- Tokyo, Japan
- Santorini, Greece

### Tour

A tour contains:

- destination
- name
- description
- price
- start date
- end date
- capacity

### Booking

A booking connects:

- user
- tour
- status

Booking statuses:

- `pending`
- `confirmed`
- `cancelled`

### Payment

A payment contains:

- booking
- amount
- status

Payment statuses:

- `pending`
- `paid`
- `failed`

### Review

A review contains:

- user
- tour
- rating
- comment

## Recommended Screens

These are the screens the designer should focus on first.

### 1. Home / Landing Page

Purpose:

- explain the product
- highlight featured destinations or tours
- give clear entry points into browsing

Suggested sections:

- hero section
- featured destinations
- featured tours
- simple value proposition
- CTA to browse tours

### 2. Destinations List

Purpose:

- show all available destinations
- help users understand where tours are available

UI should support:

- card or list layout
- empty state
- loading state
- error state

### 3. Tours List

Purpose:

- show all tours
- help users compare price, dates, and destination

Important tour info:

- tour name
- destination
- short description
- price
- date range
- capacity

### 4. Tour Detail Page

Purpose:

- present the full tour information
- allow booking
- show reviews

Important sections:

- title
- destination
- description
- schedule
- price
- capacity
- booking CTA
- reviews list
- review form

### 5. Registration Page

Purpose:

- create a new user account

Fields:

- name
- email
- password

Important note:

- there is registration, but no login flow yet

### 6. Booking Flow

Purpose:

- let a user create a booking for a tour

Current backend behavior:

- booking status is created as `pending`

Design implication:

- keep booking flow simple
- avoid multi-step checkout assumptions unless they are only conceptual

### 7. User Bookings Page

Purpose:

- show all bookings for one user
- expose booking status clearly

Important info:

- tour id or tour summary
- booking status
- creation date
- payment action if needed

### 8. Payment Flow

Purpose:

- create a payment for a booking
- show whether payment is pending, paid, or failed

Current backend behavior:

- payment is created as `pending`
- payment status can later change

Design implication:

- include clear status communication
- avoid assuming a rich checkout integration already exists

### 9. Reviews Section

Purpose:

- build trust on each tour detail page

Important info:

- rating
- comment
- date

## UX Requirements

The backend is basic, so the UI needs to do extra clarity work.

- Design strong empty states for destinations, tours, bookings, and reviews
- Design loading states for every data-driven page
- Design friendly error states
- Make status badges very clear for bookings and payments
- Favor simple user flows over advanced dashboards
- Assume lists may be short and content may be sparse

## Backend Constraints That Affect Design

These constraints are important. Do not design core flows that depend on missing backend features.

- Login is supported and returns JWT access/refresh tokens
- No logout flow exists
- No current-user profile endpoint exists
- No search endpoint exists
- No filters or sorting are implemented in the backend
- No pagination is implemented
- No favorites or wishlist feature exists
- No admin panel is implemented
- JWT authorization is enforced on protected routes
- No CORS setup is present yet, so frontend integration may use proxying in development

## Design Recommendations

### Visual Direction

The UI should feel:

- modern
- travel-focused
- clean
- optimistic
- easy to scan

Good visual ingredients:

- large destination imagery
- strong card layouts
- clear date and price presentation
- obvious call-to-action buttons
- trust-building review presentation

### Information Density

Keep the core browsing experience lightweight.

- users should be able to scan tours quickly
- details should be expanded on the tour detail page
- avoid overcrowded admin-like tables for the public-facing experience

### Mobile Behavior

The product should work well on mobile.

Prioritize:

- responsive cards
- readable forms
- clear buttons
- sticky or prominent booking CTA on tour detail

## Real Data Shape

The backend currently works with these fields.

### Destination

- `id`
- `name`
- `description`
- `image_url`
- `created_at`

### Tour

- `id`
- `destination_id`
- `name`
- `description`
- `price`
- `start_date`
- `end_date`
- `capacity`
- `created_at`

### Booking

- `id`
- `user_id`
- `tour_id`
- `status`
- `created_at`

### Payment

- `id`
- `booking_id`
- `amount`
- `status`
- `created_at`

### Review

- `id`
- `user_id`
- `tour_id`
- `rating`
- `comment`
- `created_at`

## What The Designer Can Safely Assume

- There will be a list of destinations
- There will be a list of tours
- A tour has price, schedule, and capacity
- A user can register
- A user can create a booking
- A payment exists for a booking
- Reviews belong to tours

## What The Designer Should Not Assume Yet

- Real authentication
- Social login
- Saved tours
- Rich user profile editing
- Search and filter sidebars backed by real API support
- Multi-currency support
- Promo codes
- Notification center
- Messaging/chat
- Map-based discovery backed by real geodata

## Suggested Deliverables

- desktop and mobile designs for the main screens
- reusable design system basics
- status badge styles for booking and payment states
- empty, loading, and error states
- booking and payment CTA patterns

## If You Need To Make Tradeoffs

Prefer designing:

- a strong browsing experience
- a polished tour detail page
- a simple booking journey
- clear status communication

over designing:

- complex account systems
- admin dashboards
- advanced filtering tools

## Engineering Reference

If a designer or frontend engineer needs the exact API contract, see:

- [AGENT.md](./AGENT.md)

That file contains the real routes, payload shapes, and backend limitations.
