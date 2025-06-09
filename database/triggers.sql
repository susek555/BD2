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

CREATE OR REPLACE FUNCTION delete_unsold_offers()
RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM sale_offers so 
    WHERE so.user_id = OLD.id
    AND NOT EXISTS (
        SELECT 1 FROM purchases p 
        WHERE p.offer_id = so.id
    );
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER delete_unsold_offers_trigger
BEFORE DELETE ON users
FOR EACH ROW
EXECUTE FUNCTION delete_unsold_offers();