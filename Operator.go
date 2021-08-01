package RPG_bot

import (
	"fmt"
)

type Key uint64
type Operator interface {
	InitGame(game *Game)
	StartTurn(game *Game)
	EndTurn(game *Game)
	EndGame(game *Game)
	String() string
}

type Null_Operator struct{}

func (no Null_Operator) InitGame(id Key)      {}
func (no Null_Operator) StartTurn(game *Game) {}
func (no Null_Operator) EndTurn(game *Game)   {}
func (no Null_Operator) EndGame(game *Game)   {}
func (no Null_Operator) String() string       { return "Null_Operator" }

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
