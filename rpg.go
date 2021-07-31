package rpg_bot

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
type Game struct {
	ID      Key
	Players []Player
}

var active_games map[Key]*Game

func (game *Game) RegisterPlayer(userID string) {
	for _, p := range game.Players {
		if p.id.String() == userID {
			return
		}
	}

	player := createPlayer(userID)
	game.Players = append(game.Players, player)
}

func GetGame(id Key) (*Game, error) {
	if active_games == nil {
		active_games = make(map[Key]*Game)
	}
	game, ok := active_games[id]
	if !ok {
		return nil, errors.New("game not found")
	}
	return game, nil
}

func CreateNewGame(id Key) (*Game, error) {
	game, err := GetGame(id)
	if err == nil {
		return game, errors.New("game id exists")
	}

	active_games[id] = &Game{ID: id}
	return active_games[id], nil
}

func DeleteGame(id Key) {
	_, err := GetGame(id)

	if err == nil {
		delete(active_games, id)
	}
}

// ===========================================================================
type Player struct {
	Member
	Role  string
	Score int
}

func (p *Player) AssignAndDMRole(role string) error {
	p.Role = role

	ch, err := bot_session.CreatePrivateChannel(discord.UserID(p.id))
	if err != nil {
		return err
	}

	_, err = bot_session.SendMessage(ch.ID, p.Role, nil)
	return err
}

func (p Player) String() string {
	return fmt.Sprintf("[%s] %s, %d", p.id, p.Role, p.Score)
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
	process func(e *gateway.InteractionCreateEvent)
	command api.CreateCommandData
}

var rpg_commands map[string]Command_Handler

func respondToInteraction(e *gateway.InteractionCreateEvent, response string) {
	// Respond to interaction
	data := api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content: response,
		},
	}

	if err := bot_session.RespondInteraction(e.ID, e.Token, data); err != nil {
		log.Println("failed to send interaction callback:", err)
	}
}

// TODO(Abdelrahman) Better error handling
func configureCommandHandlers(opr Operator) {
	if rpg_commands == nil {
		rpg_commands = make(map[string]Command_Handler)
	}

	rpg_commands["init-game"] = Command_Handler{
		process: func(e *gateway.InteractionCreateEvent) {
			opt := e.Data.Options
			id, _ := strconv.ParseUint(opt[0].Value, 10, 64)
			opr.InitGame(Key(id))

			response := fmt.Sprintf("Game[%d] initialized, %s!", id, opr.String())
			respondToInteraction(e, response)
		},

		// TODO(Abdelrahman) Automate id creation
		command: api.CreateCommandData{
			Name:        "init-game",
			Description: "Starts a new game, given an id",
			Options: []discord.CommandOption{
				{
					Type:        discord.StringOption,
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
			game, err := GetGame(Key(id))
			if err != nil {
				game, _ = CreateNewGame(Key(id))
			}

			userID := getUserID(e, "")
			game.RegisterPlayer(userID)

			response := fmt.Sprintf("Registered <@%s> to game[%d]!", userID, id)
			respondToInteraction(e, response)
		},

		command: api.CreateCommandData{
			Name:        "register-player",
			Description: "Register a player to a game, given an id",
			Options: []discord.CommandOption{
				{
					Type:        discord.StringOption,
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
			game, err := GetGame(Key(id))
			if err != nil {
				log.Println("failed to start turn", err)
				return
			}

			opr.StartTurn(game)

			response := fmt.Sprintf("Starting a new turn! game[%d]", id)
			respondToInteraction(e, response)
		},
		command: api.CreateCommandData{
			Name:        "start-turn",
			Description: "Starts a new turn, given a game id",
			Options: []discord.CommandOption{
				{
					Type:        discord.StringOption,
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
			game, err := GetGame(Key(id))
			if err != nil {
				log.Println("failed to end turn", err)
				return
			}

			opr.UpdateState(game)

			response := fmt.Sprintf("Ending turn! game[%d]", id)
			respondToInteraction(e, response)
		},
		command: api.CreateCommandData{
			Name:        "end-turn",
			Description: "Ends the current turn, given a game id",
			Options: []discord.CommandOption{
				{
					Type:        discord.StringOption,
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

			opr.EndGame(Key(id))

			response := fmt.Sprintf("Game Over! game[%d]", id)
			respondToInteraction(e, response)
		},
		command: api.CreateCommandData{
			Name:        "end-game",
			Description: "Ends the game, given a game id",
			Options: []discord.CommandOption{
				{
					Type:        discord.StringOption,
					Name:        "id",
					Description: "Game id",
					Required:    true,
				},
			},
		},
	}
}
