CREATE TYPE SELECTOR AS ENUM (
    'P', 'C'
);

CREATE TYPE OFFER_STATUS AS ENUM (
    'pending', 'ready', 'published', 'sold', 'expired'
);

CREATE TYPE COLOR AS ENUM (
    'red', 'blue', 'yellow', 'green', 'orange', 'purple', 'brown', 'black', 'white', 'gray',
    'cyan', 'magenta', 'lime', 'navy', 'teal', 'maroon', 'olive', 'beige', 'gold', 'other'
);

CREATE TYPE FUEL_TYPE AS ENUM (
    'diesel', 'petrol', 'electric', 'ethanol', 'lpg', 'biofuel', 'hybrid', 'hydrogen'
);

CREATE TYPE TRANSMISSION AS ENUM (
    'manual', 'automatic', 'cvt', 'dual_clutch'
);

CREATE TYPE DRIVE AS ENUM (
    'fwd', 'rwd', 'awd'
);

CREATE TYPE DOCUMENT_TYPE AS ENUM (
    'invoice', 'receipt', 'other'
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    selector SELECTOR NOT NULL
);

CREATE TABLE people (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(20) NOT NULL,
    surname VARCHAR(20) NOT NULL
);

CREATE TABLE companies (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    nip VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(30) NOT NULL
);


CREATE TABLE manufacturers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

CREATE TABLE models (
    id SERIAL PRIMARY KEY,
    manufacturer_id INTEGER NOT NULL REFERENCES manufacturers(id),
    name VARCHAR(64) NOT NULL
);

CREATE TABLE sale_offers (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    description VARCHAR(2000) NOT NULL,
    price INTEGER NOT NULL,
    date_of_issue TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    margin INTEGER NOT NULL CHECK (margin in (3, 5, 10)) ,
    status OFFER_STATUS NOT NULL,
    is_auction BOOLEAN DEFAULT FALSE
);

CREATE TABLE auctions (
    offer_id INTEGER PRIMARY KEY REFERENCES sale_offers(id) ON DELETE CASCADE,
    date_end TIMESTAMPTZ NOT NULL,
    buy_now_price INTEGER,
    initial_price INTEGER NOT NULL
);

CREATE TABLE bids (
    id SERIAL PRIMARY KEY,
    auction_id INTEGER NOT NULL REFERENCES auctions(offer_id),
    bidder_id INTEGER NOT NULL REFERENCES users(id),
    amount INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cars (
    offer_id INTEGER PRIMARY KEY REFERENCES sale_offers(id) ON DELETE CASCADE,
    vin VARCHAR(17),
    production_year INTEGER NOT NULL,
    mileage INTEGER NOT NULL,
    number_of_doors INTEGER NOT NULL CHECK (number_of_doors BETWEEN 1 AND 6),
    number_of_seats INTEGER NOT NULL CHECK (number_of_seats BETWEEN 2 AND 100),
    engine_power INTEGER NOT NULL CHECK (engine_power BETWEEN 1 AND 9999),
    engine_capacity INTEGER NOT NULL CHECK (engine_capacity BETWEEN 1 AND 9000),
    registration_number VARCHAR(20) NOT NULL,
    registration_date DATE NOT NULL,
    color COLOR NOT NULL,
    fuel_type FUEL_TYPE NOT NULL,
    drive DRIVE NOT NULL,
    transmission TRANSMISSION NOT NULL,
    number_of_gears INTEGER NOT NULL CHECK (number_of_gears BETWEEN 1 AND 10),
    model_id INTEGER REFERENCES models(id)
);

CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    offer_id INTEGER REFERENCES sale_offers(id),
    title VARCHAR(100) NOT NULL,
    description VARCHAR(200) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE client_notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    notification_id INTEGER REFERENCES notifications(id),
    seen BOOLEAN DEFAULT FALSE
);

CREATE TABLE liked_offers (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    offer_id INTEGER REFERENCES sale_offers(id)
);

CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    token VARCHAR(500) NOT NULL,
    expiry_date TIMESTAMP NOT NULL
);

CREATE TABLE purchases (
    id SERIAL PRIMARY KEY,
    offer_id INTEGER REFERENCES sale_offers(id),
    buyer_id INTEGER REFERENCES users(id),
    final_price INTEGER NOT NULL,
    issue_date DATE NOT NULL
);

CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    reviewer_id INTEGER REFERENCES users(id),
    reviewee_id INTEGER REFERENCES users(id),
    description VARCHAR(200) NOT NULL,
    rating INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
    review_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (reviewer_id <> reviewee_id),
    UNIQUE (reviewer_id, reviewee_id)
);

CREATE TABLE images (
    id SERIAL PRIMARY KEY,
    offer_id INTEGER REFERENCES sale_offers(id),
    url VARCHAR(200) NOT NULL UNIQUE,
    public_id VARCHAR(200) NOT NULL UNIQUE
);

INSERT INTO manufacturers (name) VALUES
('Toyota'),
('Volkswagen'),
('Ford'),
('BMW'),
('Mercedes-Benz'),
('Audi'),
('Honda'),
('Nissan'),
('Hyundai'),
('Kia');

INSERT INTO models (manufacturer_id, name) VALUES
(1, 'Corolla'),
(1, 'Camry'),
(2, 'Golf'),
(2, 'Passat'),
(3, 'Focus'),
(3, 'Fiesta'),
(4, '3 Series'),
(4, '5 Series'),
(5, 'C-Class'),
(5, 'E-Class'),
(6, 'A4'),
(6, 'A6'),
(7, 'Civic'),
(7, 'Accord'),
(8, 'Altima'),
(8, 'Sentra'),
(9, 'Elantra'),
(9, 'Sonata'),
(10, 'Sportage'),
(10, 'Seltos');
