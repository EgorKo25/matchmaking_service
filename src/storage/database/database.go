package database

import (
	"context"
	"errors"
	"fmt"

	"matchamking/src/config"
	"matchamking/src/core"

	"github.com/jackc/pgx/v5/pgxpool"

	_ "matchamking/migrations"
)

func NewDB(ctx context.Context, config *config.Database) (*Db, error) {
	pool, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Host,
		config.Port,
		config.DBName,
	))
	if err != nil {
		return nil, err
	}
	return &Db{pool: pool}, MigrateDatabase(pool)
}

type Db struct {
	pool *pgxpool.Pool
}

func (d *Db) Insert(ctx context.Context, user *core.Player) error {
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

func (d *Db) GetAllPlayers(ctx context.Context) ([]*core.Player, error) {
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
