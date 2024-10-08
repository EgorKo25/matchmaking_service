package core

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"matchamking/src/config"

	"github.com/stretchr/testify/require"
)

func TestCheckApproximatelyEqual(t *testing.T) {
	require.Truef(t, checkApproximatelyEqual(100, 105, 5), "expected values to be approximately equal")
	require.Falsef(t, checkApproximatelyEqual(100, 110, 5), "expected values to not be approximately equal")
}

func TestInitMatchmaker(t *testing.T) {
	// проверяем, что объект будет вызван лишь единожды
	InitMatchmaker(&config.MatchmakerConfig{AcceptableWaitingTime: 15})
	core := GetMatchmakingCore()
	InitMatchmaker(&config.MatchmakerConfig{AcceptableWaitingTime: 25})
	core1 := GetMatchmakingCore()
	require.Equal(t, core, core1)
}

func setupTestMatchmakingCore() *MatchmakingCore {
	mmCore := GetMatchmakingCore()
	if mmCore == nil {
		InitMatchmaker(
			&config.MatchmakerConfig{AcceptableWaitingTime: 5 * time.Minute, DeltaSkill: 10, DeltaLatency: 50, GroupSize: 3})
		mmCore = GetMatchmakingCore()
	}
	mmCore.GroupSize = 3
	mmCore.AcceptableWaitingTime = 5 * time.Minute
	mmCore.DeltaLatency = 50
	mmCore.DeltaSkill = 10
	return mmCore
}
func tearDownMatchmakingCore() {
	core := GetMatchmakingCore()
	core.Mutex.Lock()
	core.groups = nil
	core.Mutex.Unlock()
	matchmaker.ticker.Stop()
}

func TestFindGroup_CreateNewGroup(t *testing.T) {
	mmCore := setupTestMatchmakingCore()
	defer tearDownMatchmakingCore()
	player := randomPlayer()

	mmCore.FindGroup(player)
	require.Equal(t, 1, len(mmCore.groups))
	require.Equal(t, player, mmCore.groups[0].players[0])
	require.Equal(t, int32(1), mmCore.groups[0].totalPlayers)
}

func TestFindGroup_BoundaryConditions(t *testing.T) {
	mmCore := setupTestMatchmakingCore()
	defer tearDownMatchmakingCore()

	group := &Group{
		players:                   []*Player{},
		averagePermissibleSkill:   50.0,
		averagePermissibleLatency: 100.0,
		differenceSkill:           10.0,
		differenceLatency:         50.0,
		totalPlayers:              0,
	}
	mmCore.Mutex.Lock()
	mmCore.groups = append(mmCore.groups, group)
	mmCore.Mutex.Unlock()
	player := &Player{
		Name:    "Player4",
		Skill:   60.0,
		Latency: 150.0,
	}

	mmCore.FindGroup(player)

	require.Equal(t, 1, len(mmCore.groups))
	require.Equal(t, player, mmCore.groups[0].players[0])
}

func TestFindGroup_FillGroupAndRemove(t *testing.T) {
	mmCore := setupTestMatchmakingCore()
	defer tearDownMatchmakingCore()

	group := &Group{
		players:                   []*Player{},
		averagePermissibleSkill:   50.0,
		averagePermissibleLatency: 100.0,
		differenceSkill:           10.0,
		differenceLatency:         50.0,
		totalPlayers:              0,
	}
	mmCore.Mutex.Lock()
	mmCore.groups = append(mmCore.groups, group)
	mmCore.Mutex.Unlock()

	players := [3]*Player{
		{Name: "Player1", Skill: 55.0, Latency: 110.0},
		{Name: "Player2", Skill: 52.0, Latency: 105.0},
		{Name: "Player3", Skill: 53.0, Latency: 108.0},
	}

	for _, player := range players {
		mmCore.FindGroup(player)
	}
	require.Equal(t, 0, len(mmCore.groups))
}

// Генерация случайного игрока
func randomPlayer() *Player {
	return &Player{
		Name:    "Player",
		Skill:   rand.Float64() * 1000,
		Latency: rand.Float64() * 100,
	}
}

// BenchmarkFindGroup тестирует производительность метода FindGroup
func BenchmarkFindGroup(b *testing.B) {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Генерация 10000 случайных игроков
	players := make([]*Player, 0, 10000)
	for i := 0; i < 10000; i++ {
		players = append(players, randomPlayer())
	}

	core := &MatchmakingCore{
		groups:                []*Group{},
		GroupSize:             3,
		AcceptableWaitingTime: 5 * time.Minute,
		DeltaLatency:          10,
		DeltaSkill:            100,
	}

	// Сбрасываем таймер перед тестом
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, player := range players {
			core.FindGroup(player)
		}
	}
}

func BenchmarkFindGroupForGraph(b *testing.B) {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	file, err := os.Create("benchmark_results.txt")
	if err != nil {
		b.Fatalf("Ошибка при создании файла: %v", err)
	}
	defer file.Close()

	playerSizes := []int{1, 10, 100, 1000, 10000, 100000, 1000000}

	for _, size := range playerSizes {
		players := make([]*Player, 0, size)
		for i := 0; i < size; i++ {
			players = append(players, randomPlayer())
		}

		core := &MatchmakingCore{
			groups:                []*Group{},
			GroupSize:             3,
			AcceptableWaitingTime: 5 * time.Minute,
			DeltaLatency:          10,
			DeltaSkill:            100,
			ticker:                time.NewTicker(5 * time.Minute),
		}

		start := time.Now()

		for i := 0; i < b.N; i++ {
			for _, player := range players {
				core.FindGroup(player)
			}
		}

		duration := time.Since(start).Nanoseconds()

		_, err = fmt.Fprintf(file, "QueueSize: %d, Duration: %d\n", size, duration)
		if err != nil {
			b.Fatalf("Ошибка при записи в файл: %v", err)
		}
	}
}
