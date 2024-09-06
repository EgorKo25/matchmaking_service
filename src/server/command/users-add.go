package command

import (
	"encoding/json"
	"time"

	"matchamking/core"

	"github.com/gin-gonic/gin"
)

type UserAdd struct {
	PlayerName string  `json:"name"`
	Latency    float64 `json:"latency"`
	Skill      float64 `json:"skill"`
	CreatedAt  time.Time
}

func (u *UserAdd) Name() string {
	return "v1/user"
}
func (u *UserAdd) Parse(ctx *gin.Context) error {
	var buff []byte
	if _, err := ctx.Request.Body.Read(buff); err != nil {
		return err
	}
	if err := json.Unmarshal(buff, u); err != nil {
		return err
	}
	u.CreatedAt = time.Now()
	return nil
}

func (u *UserAdd) Apply() (any, error) {
	matchmaker := core.GetMatchmakingCore()
	matchmaker.FindGroup(u.castToMatchmakingPlayer())
	return nil, nil
}

func (u *UserAdd) castToMatchmakingPlayer() *core.Player {
	return &core.Player{
		Name:      u.Name(),
		Latency:   u.Latency,
		Skill:     u.Skill,
		CreatedAt: u.CreatedAt,
	}
}
