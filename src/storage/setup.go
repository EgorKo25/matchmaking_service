package storage

import (
	"context"
	"sync"

	"matchamking/src/config"
	"matchamking/src/core"
	"matchamking/src/storage/database"
	"matchamking/src/storage/local"
)

const (
	DatabaseStorage = iota
	LocalStorage
)

type IStorage interface {
	Insert(ctx context.Context, user *core.Player) error
	GetAllPlayers(ctx context.Context) ([]*core.Player, error)
}

var Storage IStorage
var once sync.Once

func InitStorage(ctx context.Context, config *config.Storage) error {
	var err error
	var storage IStorage
	switch config.StorageType {
	case LocalStorage:
		storage = local.NewStorage()
	case DatabaseStorage:
		storage, err = database.NewDB(ctx, config.Database)
	}
	if err != nil {
		return err
	}
	once.Do(func() {
		Storage = storage
	})
	return nil
}
