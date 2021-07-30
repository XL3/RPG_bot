package main

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
)

var bonk_command = api.CreateCommandData{
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
}
var bonk_emote string = "<:BONK:864836906898423828>"

func handleBonk(e *gateway.InteractionCreateEvent, name string) {
	userID := getUserID(e, name)
	player := createPlayer(userID)

	bonk_content := fmt.Sprintf("%s <@%s>", bonk_emote, userID)

	if err := player.assignAndDMRole(bonk_content); err != nil {
		log.Println("failed to send dm:", err)
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
}
