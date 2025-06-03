CREATE OR REPLACE FUNCTION set_is_auction_true()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE sale_offers SET is_auction = TRUE WHERE id = NEW.offer_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER auction_insert_trigger
AFTER INSERT ON auctions
FOR EACH ROW EXECUTE FUNCTION set_is_auction_true();