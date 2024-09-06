package local

import (
	"context"

	"matchamking/src/core"
)

func NewStorage() *Storage {
	return &Storage{players: make([]*core.Player, 0)}
}

type Storage struct {
	players []*core.Player
}

func (s *Storage) Insert(_ context.Context, player *core.Player) (err error) {
	s.players = append(s.players, player)
	return
}
func (s *Storage) GetAllPlayers(_ context.Context) ([]*core.Player, error) {
	return s.players, nil
}
