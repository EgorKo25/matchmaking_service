package command

import (
	"encoding/json"
	"fmt"
	"time"

	"matchamking/core"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserAdd struct {
	PlayerName string  `json:"name" validate:"required,min=2,max=100"`
	Latency    float64 `json:"latency" validate:"required,gte=0"`
	Skill      float64 `json:"skill" validate:"required,gte=0"`
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
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return fmt.Errorf("Field: %s, Error: %s\n", err.Field(), err.Tag())
		}
		return err
	}
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
