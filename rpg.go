package main

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/diamondburned/arikawa/v2/discord"
)

type Key uint64
type Operator interface {
	initSession(id Key) *Session
	startTurn(session *Session)
	updateState(session *Session)
}

// ===========================================================================
type Session struct {
	id      Key
	players []Player
}

var guild_sessions []Session

func (session *Session) registerPlayer(userID string) {
	for _, p := range session.players {
		if p.id.String() == userID {
			return
		}
	}

	player := createPlayer(userID)
	session.players = append(session.players, player)
}

func createNewSession(id Key) (*Session, error) {
	for _, session := range guild_sessions {
		if session.id == id {
			return nil, errors.New("session id exists")
		}
	}

	session := Session{}
	session.id = Key(rand.Uint64())
	return &session, nil
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
