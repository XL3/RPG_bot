package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/gateway"
)

func getUserIDAndPlayer(e *gateway.InteractionCreateEvent, name string) (string, player) {
	userID := e.Member.User.ID.String()
	if strings.HasPrefix(name, "<@!") {
		userID = name[3 : len(name)-1]
	} else {
		for _, mem := range guild_members {
			if name == mem.name || name == mem.nick {
				userID = mem.id.String()
				break
			}
		}
	}

	var ply player
	for _, mem := range guild_members {
		if userID == mem.id.String() {
			ply = player{mem, "", 0}
		}
	}

	return userID, ply
}

func handleBonk(e *gateway.InteractionCreateEvent, name string) {
	userID, ply := getUserIDAndPlayer(e, name)
	bonk_content := fmt.Sprintf("%s <@%s>", bonk_emote, userID)

	if err := ply.assignAndDMRole(bonk_content); err != nil {
		log.Println("failed to send dm:", err)
	}

	// To send a message in the channel
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
