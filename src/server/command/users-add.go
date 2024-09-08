package command

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"matchamking/src/core"
	"matchamking/src/storage"

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
	return "users"
}
func (u *UserAdd) Parse(ctx *gin.Context) error {
	buff, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(buff, u); err != nil {
		return err
	}
	u.CreatedAt = time.Now()
	validate := validator.New()
	err = validate.Struct(u)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return fmt.Errorf("field: %s, error: %s", err.Field(), err.Tag())
		}
		return err
	}
	return nil
}

func (u *UserAdd) Apply(ctx *gin.Context) (any, error) {
	matchmaker := core.GetMatchmakingCore()
	go matchmaker.FindGroup(u.castToPlayer())
	store := storage.GetStorage()
	return nil, store.Insert(ctx, u.castToPlayer())
}

func (u *UserAdd) castToPlayer() *core.Player {
	return &core.Player{
		Name:      u.PlayerName,
		Latency:   u.Latency,
		Skill:     u.Skill,
		CreatedAt: u.CreatedAt,
	}
}
