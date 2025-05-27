CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    selector VARCHAR(10) NOT NULL CHECK (selector IN ('P', 'C'))
);

CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    type VARCHAR(30) NOT NULL
);

CREATE TABLE manufacturers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
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

CREATE TABLE sale_offers (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    description VARCHAR(2000) NOT NULL,
    price INTEGER NOT NULL,
    margin FLOAT NOT NULL,
    date_of_issue TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) NOT NULL
);

CREATE TABLE auctions (
    offer_id INTEGER PRIMARY KEY REFERENCES sale_offers(id),
    date_end TIMESTAMPTZ NOT NULL,
    buy_now_price INTEGER
);


CREATE TABLE bids (
    id SERIAL PRIMARY KEY,
    auction_id INTEGER NOT NULL REFERENCES auctions(offer_id),
    bidder_id INTEGER NOT NULL REFERENCES users(id),
    amount INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE countries_of_origin (
    id SERIAL PRIMARY KEY,
    country_name VARCHAR(48) NOT NULL
);

CREATE TABLE transmissions (
    id SERIAL PRIMARY KEY,
    gear_count INTEGER NOT NULL,
    transmission_type VARCHAR(12) NOT NULL CHECK (transmission_type IN ('manual', 'automatic'))
);

CREATE TABLE colors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL
);

CREATE TABLE fuel_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL
);

CREATE TABLE drive_types (
    id SERIAL PRIMARY KEY,
    type VARCHAR(8) NOT NULL CHECK (type IN ('FWD', 'RWD', 'AWD'))
);

CREATE TABLE models (
    id SERIAL PRIMARY KEY,
    manufacturer_id INTEGER NOT NULL REFERENCES manufacturers(id),
    name VARCHAR(64) NOT NULL
);

CREATE TABLE cars (
    offer_id INTEGER PRIMARY KEY REFERENCES sale_offers(id),
    vin VARCHAR(17),
    production_year INTEGER NOT NULL,
    mileage INTEGER NOT NULL,
    number_of_doors INTEGER NOT NULL CHECK (number_of_doors BETWEEN 1 AND 6),
    number_of_seats INTEGER NOT NULL CHECK (number_of_seats BETWEEN 2 AND 100),
    engine_power INTEGER NOT NULL CHECK (engine_power BETWEEN 1 AND 9999),
    engine_capacity INTEGER NOT NULL CHECK (engine_capacity BETWEEN 1 AND 9000),
    registration_number VARCHAR(8) NOT NULL,
    registration_date DATE NOT NULL,
    color VARCHAR(20) NOT NULL,
    fuel_type VARCHAR(20) NOT NULL,
    transmission VARCHAR(20) NOT NULL,
    number_of_gears INTEGER NOT NULL CHECK (number_of_gears BETWEEN 1 AND 10),
    drive VARCHAR(8) NOT NULL CHECK (drive IN ('FWD', 'RWD', 'AWD')),
    model_id INTEGER REFERENCES models(id)
);


CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    offer_id INTEGER REFERENCES sale_offers(id),
    title VARCHAR(100) NOT NULL,
    description VARCHAR(200) NOT NULL,
    date VARCHAR(50) NOT NULL
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

CREATE TABLE photos (
    id SERIAL PRIMARY KEY,
    offer_id INTEGER REFERENCES sale_offers(id),
    url VARCHAR(50) NOT NULL
);

CREATE TABLE payment_statuses (
    id SERIAL PRIMARY KEY,
    status VARCHAR(16) NOT NULL CHECK (status IN ('pending', 'completed', 'failed', 'refunded'))
);

CREATE TABLE purchases (
    id SERIAL PRIMARY KEY,
    document_id INTEGER REFERENCES documents(id),
    payment_status_id INTEGER REFERENCES payment_statuses(id),
    offer_id INTEGER REFERENCES sale_offers(id),
    user_id INTEGER REFERENCES users(id),
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
    url VARCHAR(200) NOT NULL
)

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
