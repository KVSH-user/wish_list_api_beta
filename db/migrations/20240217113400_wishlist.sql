-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wishlist (
    wishlist_id SERIAL PRIMARY KEY ,
    name VARCHAR NOT NULL,
    uid INT NOT NULL,
    alias VARCHAR UNIQUE NOT NULL,
    FOREIGN KEY (uid) REFERENCES wish_users(uid)
);

CREATE TABLE IF NOT EXISTS items (
                                        gift_id SERIAL PRIMARY KEY ,
                                        wishlist_id INT NOT NULL,
                                        name VARCHAR NOT NULL,
                                        url VARCHAR NOT NULL,
                                        FOREIGN KEY (wishlist_id) REFERENCES wishlist(wishlist_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE items;
DROP TABLE wishlist;
-- +goose StatementEnd
