package core

import (
	"cmp"
	"fmt"
	"math"
	"slices"
	"sync"
	"time"

	"matchamking/config"
)

type Player struct {
	Name      string
	Latency   float64
	Skill     float64
	CreatedAt time.Time
}

type Group struct {
	players []*Player

	averagePermissibleLatency float64
	averagePermissibleSkill   float64

	differenceLatency float64
	differenceSkill   float64

	totalPlayers int32
	updatedAt    time.Time
}

// AddPlayer добавляет пользователя в группу
func (m *Group) AddPlayer(player *Player) {
	m.players = append(m.players, player)
	m.totalPlayers++
	m.averagePermissibleSkill = (m.averagePermissibleSkill*float64(m.totalPlayers-1) + player.Skill) / float64(m.totalPlayers)
	m.averagePermissibleLatency = (m.averagePermissibleLatency*float64(m.totalPlayers-1) + player.Latency) / float64(m.totalPlayers)
}

var matchmaker *MatchmakingCore
var once sync.Once

func GetMatchmakingCore() *MatchmakingCore {
	return matchmaker
}

// groupUpdate обновляет группы по таймеру
func (m *MatchmakingCore) groupUpdate() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _ = range m.ticker.C {
		for _, group := range m.groups {
			if time.Since(group.updatedAt) >= m.AcceptableWaitingTime {
				group.differenceLatency += m.DeltaLatency
				group.differenceSkill += m.DeltaSkill
				group.updatedAt = time.Now()
			}
		}
	}
}

func InitMatchmaker(matchmakerConfig *config.MatchmakerConfig) {
	once.Do(func() {
		matchmaker = &MatchmakingCore{
			GroupSize:             matchmakerConfig.GroupSize,
			AcceptableWaitingTime: matchmakerConfig.AcceptableWaitingTime,
			DeltaSkill:            matchmakerConfig.DeltaSkill,
			DeltaLatency:          matchmakerConfig.DeltaLatency,
			ticker:                time.NewTicker(matchmakerConfig.AcceptableWaitingTime),
		}
		go matchmaker.groupUpdate()
	})
}

type MatchmakingCore struct {
	groups []*Group

	GroupSize             int
	AcceptableWaitingTime time.Duration
	DeltaLatency          float64
	DeltaSkill            float64

	ticker *time.Ticker
	mutex  sync.Mutex
}

// formatGroupInfo выводит информацию о собранной группе
func (m *MatchmakingCore) formatGroupInfo(group *Group) {
	fmt.Printf(
		"Was create group:\nMin\\Max\\Avg latency: %0.2f\\%0.2f\\%0.2f\nMin\\Max\\Avg skill: %0.2f\\%0.2f\\%0.2f\n\n",
		slices.MinFunc(group.players, func(a, b *Player) int {
			return cmp.Compare(a.Latency, b.Latency)
		}).Latency,
		slices.MaxFunc(group.players, func(a, b *Player) int {
			return cmp.Compare(a.Latency, b.Latency)
		}).Latency,
		group.averagePermissibleLatency,
		slices.MinFunc(group.players, func(a, b *Player) int {
			return cmp.Compare(a.Skill, b.Skill)
		}).Skill,
		slices.MaxFunc(group.players, func(a, b *Player) int {
			return cmp.Compare(a.Skill, b.Skill)
		}).Skill,
		group.averagePermissibleSkill,
	)
}

// FindGroup добавляет игрока в наиболее подходящую группу или создает для него новую
func (m *MatchmakingCore) FindGroup(player *Player) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for index, group := range m.groups {
		if checkApproximatelyEqual(group.averagePermissibleSkill, player.Skill, group.differenceSkill) &&
			checkApproximatelyEqual(group.averagePermissibleLatency, player.Latency, group.differenceLatency) {
			group.AddPlayer(player)
			if len(group.players) == m.GroupSize {
				m.groups = slices.Delete(m.groups, index, index)
				go m.formatGroupInfo(group)
			}
			return
		}
	}

	m.groups = append(m.groups, &Group{
		players:                   []*Player{player},
		averagePermissibleSkill:   player.Skill,
		averagePermissibleLatency: player.Latency,
		differenceLatency:         m.DeltaLatency,
		differenceSkill:           m.DeltaSkill,
		totalPlayers:              1,
	})
}

func checkApproximatelyEqual(first, second, difference float64) bool {
	return math.Abs(first-second) <= difference
}
