package database

import (
	"github.com/jmoiron/sqlx"
)

func Migrate(db *sqlx.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS orders (
		id BIGSERIAL PRIMARY KEY,
		shop_id BIGINT NOT NULL,
		address TEXT NOT NULL,
		date TIMESTAMP NOT NULL,
		total_amount DECIMAL(10, 2) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS items (
		id BIGSERIAL PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		price DECIMAL(10, 2) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS order_items (
		id BIGSERIAL PRIMARY KEY,
		order_id BIGINT NOT NULL REFERENCES orders(id),
		item_id BIGINT NOT NULL REFERENCES items(id),
		quantity INT NOT NULL,
		total_price DECIMAL(10, 2) NOT NULL
	);
	`
	_, err := db.Exec(query)
	return err
}
