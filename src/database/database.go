package database

import (
	"context"
	"errors"

	"matchamking/src/core"

	"github.com/jackc/pgx/v5/pgxpool"

	_ "matchamking/migrations"
)

func NewDB(ctx context.Context, connString string) (*DB, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	return &DB{pool: pool}, MigrateDatabase(pool)
}

type DB struct {
	pool *pgxpool.Pool
}

func (d *DB) Insert(ctx context.Context, user *core.Player) error {
	exec, err := d.pool.Exec(ctx, `INSERT INTO players
    (name, latency, skill, created_at) VALUES ($1, $2, $3, $4)`,
		user.Name, user.Latency, user.Skill, user.CreatedAt)
	if err != nil {
		return err
	}
	if exec.RowsAffected() == 0 {
		return errors.New("player was not insert into database ")
	}
	return err
}

func (d *DB) GetAllPlayers(ctx context.Context) ([]*core.Player, error) {
	rows, err := d.pool.Query(ctx, `SELECT * FROM players`)
	if err != nil {
		return nil, err
	}
	var players []*core.Player
	for rows.Next() {
		var player core.Player
		err = rows.Scan(&player.Name, &player.Skill, &player.Latency, &player.CreatedAt)
		if err != nil {
			return nil, err
		}
		players = append(players, &player)
	}
	return players, nil
}
