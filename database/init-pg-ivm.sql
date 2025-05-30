CREATE EXTENSION IF NOT EXISTS pg_ivm;

GRANT USAGE ON SCHEMA public TO bd2_user;
GRANT CREATE ON SCHEMA public TO bd2_user;


SELECT pgivm.create_immv(
  'offer_creators',          
  $$ SELECT o.user_id  AS user_id,
           o.id       AS offer_id
     FROM  sale_offers o $$);

CREATE UNIQUE INDEX ON offer_creators (user_id, offer_id);

SELECT pgivm.create_immv(
  'offer_bidders',
  $$ SELECT b.bidder_id  AS user_id,
           b.auction_id AS offer_id
     FROM  bids b
     JOIN  sale_offers o ON o.id = b.auction_id $$);

CREATE UNIQUE INDEX ON offer_bidders (user_id, offer_id);


SELECT pgivm.create_immv(
  'offer_likers',
  $$ SELECT l.user_id  AS user_id,
           l.offer_id AS offer_id
     FROM  liked_offers l $$);

CREATE UNIQUE INDEX ON offer_likers (user_id, offer_id);


CREATE VIEW user_offer_interactions AS
SELECT user_id, offer_id FROM offer_creators
UNION
SELECT user_id, offer_id FROM offer_bidders
UNION
SELECT user_id, offer_id FROM offer_likers;

