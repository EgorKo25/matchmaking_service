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

var storage IStorage
var once sync.Once

func GetStorage() IStorage {
	return storage
}

func InitStorage(ctx context.Context, config *config.Storage) error {
	var err error
	var s IStorage
	switch config.StorageType {
	case LocalStorage:
		s = local.NewStorage()
	case DatabaseStorage:
		s, err = database.NewDB(ctx, config.Database)
	}
	if err != nil {
		return err
	}
	once.Do(func() {
		storage = s
	})
	return nil
}
