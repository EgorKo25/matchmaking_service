package command

import (
	"matchamking/src/core"
	"matchamking/src/storage"

	"github.com/gin-gonic/gin"
)

type UsersList struct {
	Users []*Player `json:"players"`
}

type Player struct {
	PlayerName string  `json:"name" validate:"required,min=2,max=100"`
	Latency    float64 `json:"latency" validate:"required,gte=0"`
	Skill      float64 `json:"skill" validate:"required,gte=0"`
}

func (u *UsersList) Name() string {
	return "users/list"
}

func (u *UsersList) Parse(_ *gin.Context) error {
	return nil
}

func (u *UsersList) Apply(ctx *gin.Context) (any, error) {
	mmStorage := storage.GetStorage()
	players, err := mmStorage.GetAllPlayers(ctx)
	if err != nil {
		return nil, err
	}
	u.Users = custToListPlayers(players)
	return u, err
}

func custToListPlayers(players []*core.Player) []*Player {
	result := make([]*Player, 0)
	for _, v := range players {
		result = append(result, &Player{
			PlayerName: v.Name,
			Latency:    v.Latency,
			Skill:      v.Skill,
		})
	}
	return result
}
