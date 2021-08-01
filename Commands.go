package RPG_bot

import (
	"fmt"
	"log"
	"strconv"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
)

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

			// Assign channel to game
			game := opr.InitGame(Key(id))
			game.ChID = e.ChannelID

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

	rpg_commands["bonk"] = Command_Handler{
		process: func(e *gateway.InteractionCreateEvent) {
			opt := e.Data.Options
			name := opt[0].Value
			userID := getUserID(e, name)
			player := createPlayer(userID)

			bonk_content := fmt.Sprintf("%s %s", bonk_emote, player)
			game := Game{}
			if err := game.MessagePlayer(player, bonk_content); err != nil {
				log.Println("failed to send dm:", err)
				return
			}

			// Respond to interaction
			data := api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Content: "Bonked!",
				},
			}

			if err := bot_session.RespondInteraction(e.ID, e.Token, data); err != nil {
				log.Println("failed to send interaction callback:", err)
			}
		},
		command: api.CreateCommandData{
			Name:        "bonk",
			Description: "Bonks a given user",
			Options: []discord.CommandOption{
				{
					Type:        discord.StringOption,
					Name:        "whom",
					Description: "User to bonk",
					Required:    true,
				},
			},
		},
	}
}

var bonk_emote string = "<:BONK:864836906898423828>"
