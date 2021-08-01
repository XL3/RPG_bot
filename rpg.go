package RPG_bot

import (
	"fmt"
)

type Key uint64
type Operator interface {
	InitGame(id Key) *Game
	StartTurn(game *Game)
	UpdateState(game *Game)
	EndGame(id Key)
	String() string
}

type Null_Operator struct{}

func (no Null_Operator) InitGame(id Key) *Game  { return nil }
func (no Null_Operator) StartTurn(game *Game)   {}
func (no Null_Operator) UpdateState(game *Game) {}
func (no Null_Operator) EndGame(id Key)         {}
func (no Null_Operator) String() string         { return "Null" }

// ===========================================================================
type Player struct {
	Member
	Role  string
	Score int
}

// Mention player
func (p Player) String() string {
	return fmt.Sprintf("<@%s>", p.id)
}

func createPlayer(userID string) Player {
	var p Player
	for _, mem := range guild_members {
		if userID == mem.id.String() {
			p = Player{mem, "", 0}
		}
	}
	return p
}
