CREATE TABLE IF NOT EXISTS users (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    email      VARCHAR(100) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    role       VARCHAR(20)  NOT NULL DEFAULT 'client',
    created_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS destinations (
    id                SERIAL PRIMARY KEY,
    name              VARCHAR(100) NOT NULL UNIQUE,
    description       TEXT,
    image_url         TEXT,
    created_at        TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tours (
    id             SERIAL PRIMARY KEY,
    destination_id INT            NOT NULL REFERENCES destinations(id),
    name           VARCHAR(200)   NOT NULL,
    description    TEXT,
    price          DECIMAL(10,2)  NOT NULL,
    start_date     TIMESTAMP      NOT NULL,
    end_date       TIMESTAMP      NOT NULL,
    capacity       INT            NOT NULL,
    created_at     TIMESTAMP      NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS bookings (
    id         SERIAL PRIMARY KEY,
    user_id    INT          NOT NULL REFERENCES users(id),
    tour_id    INT          NOT NULL REFERENCES tours(id),
    status     VARCHAR(20)  NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS payments (
    id         SERIAL PRIMARY KEY,
    booking_id INT           NOT NULL REFERENCES bookings(id),
    amount     DECIMAL(10,2) NOT NULL,
    status     VARCHAR(20)   NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP     NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS reviews (
    id         SERIAL PRIMARY KEY,
    user_id    INT       NOT NULL REFERENCES users(id),
    tour_id    INT       NOT NULL REFERENCES tours(id),
    rating     INT       NOT NULL,
    comment    TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
 );

INSERT INTO users (name, email, password, role)
VALUES ('Администратор', 'admin@tourism.tj', 'admin123', 'admin');
