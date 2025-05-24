CREATE TABLE client (
    id SERIAL PRIMARY KEY,
    login VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(50) NOT NULL,
    address VARCHAR(50) NOT NULL,
    client_type VARCHAR(10) NOT NULL CHECK (client_type IN ('individual', 'company'))
);

CREATE TABLE document (
    id SERIAL PRIMARY KEY,
    type VARCHAR(30) NOT NULL
);

CREATE TABLE manufacturer (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

CREATE TABLE individual (
    id INTEGER PRIMARY KEY REFERENCES client(id),
    first_name VARCHAR(20) NOT NULL,
    last_name VARCHAR(20) NOT NULL
);

CREATE TABLE company (
    id INTEGER PRIMARY KEY REFERENCES client(id),
    nip INTEGER NOT NULL,
    name VARCHAR(30) NOT NULL
);

CREATE TABLE listing (
    id SERIAL PRIMARY KEY,
    client_id INTEGER NOT NULL REFERENCES client(id),
    description VARCHAR(2000) NOT NULL,
    address VARCHAR(50) NOT NULL,
    price INTEGER NOT NULL,
    margin FLOAT NOT NULL,
    is_active CHAR(1) NOT NULL CHECK (is_active IN ('Y', 'N')),
    listing_type VARCHAR(10) NOT NULL CHECK (listing_type IN ('auction', 'standard'))  -- CHANGED: ad_type -> listing_type
);

CREATE TABLE auction (
    id INTEGER PRIMARY KEY REFERENCES listing(id),
    end_date DATE NOT NULL,
    buy_now_price INTEGER,
    CHECK (listing_type = 'auction')
) INHERITS (listing);

CREATE TABLE standard_offer (
    id INTEGER PRIMARY KEY REFERENCES listing(id),
    CHECK (listing_type = 'standard')
) INHERITS (listing);

CREATE TABLE bid (
    id SERIAL PRIMARY KEY,
    auction_id INTEGER NOT NULL REFERENCES auction(id),
    client_id INTEGER NOT NULL REFERENCES client(id),
    amount INTEGER NOT NULL,
    time TIMESTAMP NOT NULL
);

CREATE TABLE country_of_origin (
    id SERIAL PRIMARY KEY,
    country_name VARCHAR(48) NOT NULL
);

CREATE TABLE transmission (
    id SERIAL PRIMARY KEY,
    gear_count INTEGER NOT NULL,
    transmission_type VARCHAR(12) NOT NULL CHECK (transmission_type IN ('manual', 'automatic'))
);

CREATE TABLE color (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL
);

CREATE TABLE fuel_type (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL
);

CREATE TABLE drive_type (
    id SERIAL PRIMARY KEY,
    type VARCHAR(8) NOT NULL CHECK (type IN ('FWD', 'RWD', 'AWD'))
);

CREATE TABLE car_model (
    id SERIAL PRIMARY KEY,
    manufacturer_id INTEGER NOT NULL REFERENCES manufacturer(id),
    name VARCHAR(64) NOT NULL
);

CREATE TABLE car (
    id SERIAL PRIMARY KEY,
    production_year INTEGER NOT NULL,
    mileage INTEGER NOT NULL,
    door_count INTEGER NOT NULL,
    seat_count INTEGER NOT NULL,
    power INTEGER NOT NULL,
    first_registration_date DATE NOT NULL,
    registration_number VARCHAR(8) NOT NULL,
    engine_capacity FLOAT NOT NULL,
    fuel_type_id INTEGER REFERENCES fuel_type(id),
    transmission_id INTEGER REFERENCES transmission(id),
    color_id INTEGER REFERENCES color(id),
    drive_type_id INTEGER REFERENCES drive_type(id),
    country_of_origin_id INTEGER REFERENCES country_of_origin(id),
    model_id INTEGER REFERENCES car_model(id)
);


CREATE TABLE notification (
    id SERIAL PRIMARY KEY,
    listing_id INTEGER REFERENCES listing(id),
    title VARCHAR(100) NOT NULL,
    content VARCHAR(200) NOT NULL
);

CREATE TABLE client_notification (
    id SERIAL PRIMARY KEY,
    client_id INTEGER REFERENCES client(id),
    notification_id INTEGER REFERENCES notification(id)
);

CREATE TABLE favorite_listing (
    id SERIAL PRIMARY KEY,
    client_id INTEGER REFERENCES client(id),
    listing_id INTEGER REFERENCES listing(id)
);

CREATE TABLE refresh_token (
    client_id INTEGER PRIMARY KEY REFERENCES client(id),
    token VARCHAR(64) NOT NULL
);

CREATE TABLE photo (
    listing_id INTEGER REFERENCES listing(id),
    id SERIAL PRIMARY KEY,
    url VARCHAR(50) NOT NULL
);

CREATE TABLE payment_status (
    id SERIAL PRIMARY KEY,
    status VARCHAR(16) NOT NULL CHECK (status IN ('pending', 'completed', 'failed', 'refunded'))
);

CREATE TABLE purchase (
    id SERIAL PRIMARY KEY,
    document_id INTEGER REFERENCES document(id),
    payment_status_id INTEGER REFERENCES payment_status(id),
    listing_id INTEGER REFERENCES listing(id),
    client_id INTEGER REFERENCES client(id),
    issue_date DATE NOT NULL
);

CREATE TABLE review (
    id SERIAL PRIMARY KEY,
    reviewing_client_id INTEGER REFERENCES client(id),
    reviewed_client_id INTEGER REFERENCES client(id),
    description VARCHAR(200) NOT NULL,
    rating INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5)
);

ALTER TABLE car ADD CONSTRAINT car_transmission_fk FOREIGN KEY (transmission_id) REFERENCES transmission(id);
ALTER TABLE car ADD CONSTRAINT car_color_fk FOREIGN KEY (color_id) REFERENCES color(id);
ALTER TABLE car ADD CONSTRAINT car_drive_type_fk FOREIGN KEY (drive_type_id) REFERENCES drive_type(id);
ALTER TABLE car ADD CONSTRAINT car_country_fk FOREIGN KEY (country_of_origin_id) REFERENCES country_of_origin(id);
ALTER TABLE car ADD CONSTRAINT car_model_fk FOREIGN KEY (model_id) REFERENCES car_model(id);

ALTER TABLE car_model ADD CONSTRAINT model_manufacturer_fk FOREIGN KEY (manufacturer_id) REFERENCES manufacturer(id);