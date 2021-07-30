package main

import (
	"log"
	"math/rand"

	"github.com/diamondburned/arikawa/v2/gateway"
)

type Roulette struct {
}

func (r Roulette) initSession(id Key) *Session {
	session, err := createNewSession(id)

	if err != nil {
		log.Fatal("Failed to initialize session", err)
	}

	return session
}

func (r Roulette) startTurn(session *Session) {
	size := len(session.players)
	target := rand.Intn(size)

	for i, player := range session.players {
		if i != target {
			player.assignAndDMRole("Spared")
		} else {
			// Shoot the target
			player.assignAndDMRole("Shot")
			player.score = -1
		}

		log.Println(session.players[i])
	}
}

func (r Roulette) updateState(session *Session) {
	size := len(session.players)
	for i, player := range session.players {
		if player.score < 0 {
			session.players[i] = session.players[size-1]
			session.players = session.players[:size-1]
			return
		}
	}
}

func Roulette_registerPlayer(e *gateway.InteractionCreateEvent, session *Session) {
	userID := e.Member.User.ID.String()
	session.registerPlayer(userID)
}
