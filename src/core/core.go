package core

import (
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
}

// AddPlayer добавляет пользователя в группу
func (m *Group) AddPlayer(player *Player) {
	m.players = append(m.players, player)
	m.totalPlayers++
	m.averagePermissibleSkill = m.averagePermissibleSkill * (float64(m.totalPlayers+1) + player.Skill) / float64(m.totalPlayers)
	m.averagePermissibleLatency = m.averagePermissibleLatency * (float64(m.totalPlayers-1) + player.Latency) / float64(m.totalPlayers)
}

var matchmaker *MatchmakingCore
var once sync.Once

func GetMatchmakingCore() *MatchmakingCore {
	return matchmaker
}

func InitMatchmaker(matchmakerConfig *config.MatchmakerConfig) {
	once.Do(func() {
		matchmaker = &MatchmakingCore{
			GroupSize:             matchmakerConfig.GroupSize,
			AcceptableWaitingTime: matchmakerConfig.AcceptableWaitingTime,
			DeltaSkill:            matchmakerConfig.DeltaSkill,
			DeltaLatency:          matchmakerConfig.DeltaLatency,
		}
	})
}

type MatchmakingCore struct {
	groups []*Group

	GroupSize             int
	AcceptableWaitingTime time.Duration
	DeltaLatency          float64
	DeltaSkill            float64
}

// formatGroupInfo выводит информацию о собранной группе
func (m *MatchmakingCore) formatGroupInfo(group *Group) {
	fmt.Printf("was create group: %v\nWith average latency: %0.2f\nWith average skill: %0.2f\n",
		group.players, group.averagePermissibleLatency, group.averagePermissibleSkill)
}

// FindGroup добавляет игрока в наиболее подходящую группу или создает для него новую
func (m *MatchmakingCore) FindGroup(player *Player) {
	for index, group := range m.groups {
		if checkApproximatelyEqual(group.averagePermissibleSkill, player.Skill, group.differenceSkill) &&
			checkApproximatelyEqual(group.averagePermissibleLatency, player.Latency, group.differenceLatency) {
			group.AddPlayer(player)
			if len(group.players) == m.GroupSize {
				m.groups = slices.Delete(m.groups, index, index)
				m.formatGroupInfo(group)
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
	})
}

func checkApproximatelyEqual(first, second, difference float64) bool {
	return math.Abs(first-second) <= difference
}
