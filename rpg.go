package main

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
)

type player struct {
	member
	role  string
	score int
}

func (p *player) assignAndDMRole(role string) error {
	p.role = role

	ch, err := bot_session.CreatePrivateChannel(discord.UserID(p.id))
	if err != nil {
		return err
	}

	_, err = bot_session.SendMessage(ch.ID, p.role, nil)
	return err
}

func (p player) String() string {
	return fmt.Sprintf("[%s] %s, %d", p.id, p.role, p.score)
}
