CREATE EXTENSION IF NOT EXISTS pg_ivm;

GRANT USAGE ON SCHEMA public TO bd2_user;
GRANT CREATE ON SCHEMA public TO bd2_user;


SELECT pgivm.create_immv(
  'offer_creators',
  $$ SELECT o.user_id  AS user_id,
           o.id       AS offer_id
     FROM  sale_offers o $$
);

CREATE UNIQUE INDEX ON offer_creators (user_id, offer_id);

SELECT pgivm.create_immv(
  'offer_bidders',
  $$ SELECT DISTINCT
           b.bidder_id  AS user_id,
           b.auction_id AS offer_id
     FROM  bids b
     JOIN  sale_offers o ON o.id = b.auction_id WHERE o.status = 'published' $$
);

SELECT pgivm.create_immv(
  'offer_likers',
  $$ SELECT l.user_id  AS user_id,
           l.offer_id AS offer_id
     FROM  liked_offers l
     JOIN  sale_offers o ON o.id = l.offer_id WHERE o.status = 'published' $$
);

CREATE UNIQUE INDEX ON offer_likers (user_id, offer_id);


CREATE VIEW user_offer_interactions AS
SELECT user_id, offer_id FROM offer_creators
UNION
SELECT user_id, offer_id FROM offer_bidders
UNION
SELECT user_id, offer_id FROM offer_likers;

SELECT pgivm.create_immv(
  'sale_offer_view',
  $$ SELECT
    s.id,
    s.user_id,
    u.username,
    s.description,
    s.price,
    s.date_of_issue,
    s.margin,
    s.status,
    c.vin,
    c.production_year,
    c.mileage,
    c.number_of_doors,
    c.number_of_seats,
    c.engine_power,
    c.engine_capacity,
    c.registration_number,
    c.registration_date,
    c.color,
    c.fuel_type,
    c.transmission,
    c.number_of_gears,
    c.drive,
    man.name as brand,
    mod.name as model
    FROM sale_offers s
    JOIN users u ON u.id = s.user_id
    JOIN cars c ON c.offer_id = s.id
    JOIN models mod ON c.model_id = mod.id
    JOIN manufacturers man ON mod.manufacturer_id = man.id $$
);

CREATE UNIQUE INDEX ON sale_offer_view (id);