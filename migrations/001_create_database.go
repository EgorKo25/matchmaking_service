package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateDatabase, downCreateDatabase)
}

func upCreateDatabase(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE players (
			id SERIAL PRIMARY KEY,
    		name VARCHAR(100),
    		skill DOUBLE PRECISION,
    		latency DOUBLE PRECISION,
    		created_at TIMESTAMP
		);
	`)
	return err
}

func downCreateDatabase(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.ExecContext(ctx, `DROP TABLE users`)
	return err
}
