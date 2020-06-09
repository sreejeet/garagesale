package schema

import "github.com/jmoiron/sqlx"

const seed = `INSERT INTO products (product_id, name, cost, quantity, date_created, date_updated) VALUES
	('a2b0639f-2cc6-44b8-b97b-15d69dbb511e', 'Comic Books', 50, 42, '2019-01-01 00:00:01.000001+00', '2019-01-01 00:00:01.000001+00'),
	('72f8b983-3eb4-48db-9ed0-e45cc6bd716b', 'McDonalds Toys', 75, 120, '2019-01-01 00:00:02.000001+00', '2019-01-01 00:00:02.000001+00')
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
