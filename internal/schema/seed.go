package schema

import "github.com/jmoiron/sqlx"

const seed = `
	
	-- Create sample products
	INSERT INTO products (product_id, name, cost, quantity, date_created, date_updated) VALUES
	('a2b0639f-2cc6-44b8-b97b-15d69dbb511e', 'Comic Books', 50, 42, '2019-01-01 00:00:01.000001+00', '2019-01-01 00:00:01.000001+00'),
	('72f8b983-3eb4-48db-9ed0-e45cc6bd716b', 'McDonalds Toys', 75, 120, '2019-01-01 00:00:02.000001+00', '2019-01-01 00:00:02.000001+00')
	ON CONFLICT DO NOTHING;

	-- Create sample sales
	INSERT INTO sales (sale_id, product_id, quantity, paid, date_created) VALUES
	('98b6d4b8-f04b-4c79-8c2e-a0aef46854b7', 'a2b0639f-2cc6-44b8-b97b-15d69dbb511e', 2, 100, '2019-01-01 00:00:03.000001+00'),
	('85f6fb09-eb05-4874-ae39-82d1a30fe0d7', 'a2b0639f-2cc6-44b8-b97b-15d69dbb511e', 5, 250, '2019-01-01 00:00:04.000001+00'),
	('a235be9e-ab5d-44e6-a987-fa1c749264c7', '72f8b983-3eb4-48db-9ed0-e45cc6bd716b', 3, 225, '2019-01-01 00:00:05.000001+00')
	ON CONFLICT DO NOTHING;

	-- Create admin and regular users with password "gophers"
	INSERT INTO users (user_id, name, email, roles, password_hash, date_created, date_updated) VALUES
	('5cf37266-3473-4006-984f-9325122678b7', 'Admin Gopher', 'admin@example.com', '{ADMIN,USER}', '$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'User Gopher', 'user@example.com', '{USER}', '$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
	ON CONFLICT DO NOTHING;`

// Seed runs the above query to add some data and bring the database into a usefule state.
func Seed(db *sqlx.DB) error {

	// Using transactions in case a rollback is needed when errors occur
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(seed); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}
