package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
)

type Key uint64
type Operator interface {
	initGame(id Key) *Game
	startTurn(game *Game)
	updateState(game *Game)
	endGame(id Key)
}

type Null_Operator struct{}

func (no Null_Operator) initGame(id Key) *Game  { return nil }
func (no Null_Operator) startTurn(game *Game)   {}
func (no Null_Operator) updateState(game *Game) {}
func (no Null_Operator) endGame(id Key) {}

// ===========================================================================
type Game struct {
	id      Key
	players []Player
}

var active_games map[Key]*Game

func (game *Game) registerPlayer(userID string) {
	for _, p := range game.players {
		if p.id.String() == userID {
			return
		}
	}

	player := createPlayer(userID)
	game.players = append(game.players, player)
}

func getGame(id Key) (*Game, error) {
	game, ok := active_games[id]
	if !ok {
		return nil, errors.New("game not found")
	}
	return game, nil
}

func createNewGame(id Key) (*Game, error) {
	game, err := getGame(id)
	if err == nil {
		return game, errors.New("game id exists")
	}

	active_games[id] = &Game{id: id}
	return active_games[id], nil
}


// ===========================================================================
type Player struct {
	Member
	role  string
	score int
}

func (p *Player) assignAndDMRole(role string) error {
	p.role = role

	ch, err := bot_session.CreatePrivateChannel(discord.UserID(p.id))
	if err != nil {
		return err
	}

	_, err = bot_session.SendMessage(ch.ID, p.role, nil)
	return err
}

func (p Player) String() string {
	return fmt.Sprintf("[%s] %s, %d", p.id, p.role, p.score)
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

// ===========================================================================
// Commands
type Command_Handler struct {
	process func (e *gateway.InteractionCreateEvent)
	command api.CreateCommandData
}

var rpg_commands map[string]Command_Handler
func configureCommandHandlers(opr Operator) {
	rpg_commands["init-game"] = Command_Handler{
		process: func(e *gateway.InteractionCreateEvent) {
			opt := e.Data.Options
			id, _ := strconv.ParseUint(opt[0].Value, 10, 64)
			opr.initGame(Key(id))
		},

		// TODO(Abdelrahman) Automate id creation
		command: api.CreateCommandData{
			Name:        "init-game",
			Description: "Starts a new game, given an id",
			Options:     []discord.CommandOption{
				{
					Type: discord.StringOption,
					Name:        "id",
					Description: "Game id",
					Required:    true,
				},
			},
		},
	}

	rpg_commands["register-player"] = Command_Handler{
		process: func(e *gateway.InteractionCreateEvent) {
			opt := e.Data.Options
			id, _ := strconv.ParseUint(opt[0].Value, 10, 64)
			game, err := getGame(Key(id))
			if err != nil {
				game, _ = createNewGame(Key(id))
			}

			userID := getUserID(e, "")
			game.registerPlayer(userID)
		},

		command: api.CreateCommandData{
			Name:        "register-player",
			Description: "Register a player to a game, given an id",
			Options:     []discord.CommandOption{
				{
					Type: discord.StringOption,
					Name:        "id",
					Description: "Game id",
					Required:    true,
				},
			},
		},
	}

	rpg_commands["start-turn"] = Command_Handler{
		process: func(e *gateway.InteractionCreateEvent) {
			opt := e.Data.Options
			id, _ := strconv.ParseUint(opt[0].Value, 10, 64)
			game, err := getGame(Key(id))
			if err != nil {
				log.Println("failed to start turn", err)
				return
			}

			opr.startTurn(game)
		},
		command: api.CreateCommandData{
			Name:        "start-turn",
			Description: "Starts a new turn, given a game id",
			Options:     []discord.CommandOption{
				{
					Type: discord.StringOption,
					Name:        "id",
					Description: "Game id",
					Required:    true,
				},
			},
		},
	}
	rpg_commands["end-turn"] = Command_Handler{
		process: func(e *gateway.InteractionCreateEvent) {
			opt := e.Data.Options
			id, _ := strconv.ParseUint(opt[0].Value, 10, 64)
			game, err := getGame(Key(id))
			if err != nil {
				log.Println("failed to end turn", err)
				return
			}

			opr.updateState(game)
		},
		command: api.CreateCommandData{
			Name:        "end-turn",
			Description: "Ends the current turn, given a game id",
			Options:     []discord.CommandOption{
				{
					Type: discord.StringOption,
					Name:        "id",
					Description: "Game id",
					Required:    true,
				},
			},
		},
	}

	rpg_commands["end-game"] = Command_Handler{
		process: func(e *gateway.InteractionCreateEvent) {
			opt := e.Data.Options
			id, _ := strconv.ParseUint(opt[0].Value, 10, 64)

			opr.endGame(Key(id))
		},
		command: api.CreateCommandData{
			Name:        "end-game",
			Description: "Ends the game, given a game id",
			Options:     []discord.CommandOption{
				{
					Type: discord.StringOption,
					Name:        "id",
					Description: "Game id",
					Required:    true,
				},
			},
		},
	}
}
