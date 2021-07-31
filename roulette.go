package main

import (
	"log"
	"math/rand"
)

type Roulette struct{}

func (r Roulette) initGame(id Key) *Game {
	session, err := createNewGame(id)

	if err != nil {
		log.Fatal("Failed to initialize session", err)
	}

	return session
}

func (r Roulette) startTurn(game *Game) {
	size := len(game.players)
	target := rand.Intn(size)

	for i, player := range game.players {
		if i != target {
			player.assignAndDMRole("Spared")
		} else {
			// Shoot the target
			player.assignAndDMRole("Shot")
			player.score = -1
		}

		log.Println(game.players[i])
	}
}

func (r Roulette) updateState(game *Game) {
	size := len(game.players)
	for i, player := range game.players {
		if player.score < 0 {
			game.players[i] = game.players[size-1]
			game.players = game.players[:size-1]
			return
		}
	}
}

func (r Roulette) endGame(id Key) {
	_, err := getGame(id)

	if err == nil {
		delete(active_games, id)
	}
}
