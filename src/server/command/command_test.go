package command

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"matchamking/src/config"
	"matchamking/src/core"
	"matchamking/src/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// Mock storage to simulate Insert function
type MockStorage struct {
	players []*core.Player
}

func (ms *MockStorage) Insert(_ context.Context, player *core.Player) error {
	ms.players = append(ms.players, player)
	return nil
}

func (ms *MockStorage) GetAllPlayers(_ context.Context) ([]*core.Player, error) {
	return ms.players, nil
}

func setupTestMatchmakingCore() *core.MatchmakingCore {
	mmCore := core.GetMatchmakingCore()
	if mmCore == nil {
		core.InitMatchmaker(
			&config.MatchmakerConfig{AcceptableWaitingTime: 5 * time.Minute, DeltaSkill: 10, DeltaLatency: 50, GroupSize: 3})
		mmCore = core.GetMatchmakingCore()
	}
	mmCore.GroupSize = 3
	mmCore.AcceptableWaitingTime = 5 * time.Minute
	mmCore.DeltaLatency = 50
	mmCore.DeltaSkill = 10
	return mmCore
}

func TestUserAdd_ValidRequest(t *testing.T) {
	setupTestMatchmakingCore()
	gin.SetMode(gin.TestMode)
	mockStorage := &MockStorage{}
	storage.Storage = mockStorage

	user := UserAdd{
		PlayerName: "testuser",
		Latency:    30.5,
		Skill:      70.0,
	}
	body, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	userAddCommand := &UserAdd{}
	err := userAddCommand.Parse(ctx)
	require.NoError(t, err)

	_, err = userAddCommand.Apply(ctx)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, 1, len(mockStorage.players))
	require.Equal(t, "testuser", mockStorage.players[0].Name)
	require.InDelta(t, 30.5, mockStorage.players[0].Latency, 0.1)
	require.InDelta(t, 70.0, mockStorage.players[0].Skill, 0.1)
}

func TestUserAdd_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{"name": "", "latency": -10, "skill": -5}`)
	req, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	userAddCommand := &UserAdd{}
	err := userAddCommand.Parse(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "field: PlayerName")

	body = []byte(`{"name": "testuser", "latency": -10, "skill": -5}`)
	req, _ = http.NewRequest("POST", "/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = req

	err = userAddCommand.Parse(ctx)
	require.Contains(t, err.Error(), "field: Latency")

	body = []byte(`{"name": "testuser", "latency": 10, "skill": -5}`)
	req, _ = http.NewRequest("POST", "/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = req

	err = userAddCommand.Parse(ctx)
	require.Contains(t, err.Error(), "field: Skill")
}
func TestUsersList_ValidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockStorage := &MockStorage{
		players: []*core.Player{
			{Name: "Player1", Latency: 30.5, Skill: 70.0},
			{Name: "Player2", Latency: 40.5, Skill: 80.0},
		},
	}
	storage.Storage = mockStorage

	req, _ := http.NewRequest("POST", "/v1/users/list", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	usersListCommand := &UsersList{}
	body, err := usersListCommand.Apply(ctx)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, w.Code)

	response, ok := body.(*UsersList)
	require.True(t, ok)

	require.Equal(t, 2, len(response.Users))
	require.Equal(t, "Player1", response.Users[0].PlayerName)
	require.Equal(t, 30.5, response.Users[0].Latency)
	require.Equal(t, 70.0, response.Users[0].Skill)

	require.Equal(t, "Player2", response.Users[1].PlayerName)
	require.Equal(t, 40.5, response.Users[1].Latency)
	require.Equal(t, 80.0, response.Users[1].Skill)
}
