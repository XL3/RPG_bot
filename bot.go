package main

import (
	"log"
	"fmt"
	"strings"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/api"
)

func handle_bonk(e *gateway.InteractionCreateEvent, name string) {
	log.Println(name)
	userID := e.Member.User.ID.String()

	if strings.HasPrefix(name, "<@!") {
		userID = name[3:len(name)-1]
	} else {
		for _, mem := range guild_members {
			if name == mem.name || name == mem.nick {
				userID = mem.ID.String()
				break
			}
		}
	}
	bonk_content := fmt.Sprintf("%s <@%s>", bonk_emote, userID)

	data := api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content: bonk_content,
		},
	}

	if err := bot_session.RespondInteraction(e.ID, e.Token, data); err != nil {
		log.Println("failed to send interaction callback:", err)
	}
}
