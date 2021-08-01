package main

import (
	"fmt"
	"log"
	"math/rand"

	rb "github.com/XL3/RPG_bot"
)

type Roulette struct{}

func (r Roulette) InitGame(id rb.Key) *rb.Game {
	game, err := rb.CreateNewGame(id)
	if err != nil {
		log.Fatal("Failed to initialize game", err)
	}

	return game
}

func (r Roulette) StartTurn(game *rb.Game) {
	size := len(game.Players)
	target := rand.Intn(size)

	var msg string
	for i, player := range game.Players {
		if i != target {
			player.Role = "spared"
			msg = "You're in luck. You've been spared."
		} else {
			player.Role = "shot"
			msg = "You drew the short straw. You've been shot."
			player.Score = -1
		}

		go game.MessagePlayer(*player, msg)
	}
}

func (r Roulette) EndTurn(game *rb.Game) {
	size := len(game.Players)
	for i, player := range game.Players {
		// Remove the current player
		if player.Score < 0 {
			game.Players[i] = game.Players[size-1]
			game.Players = game.Players[:size-1]

			msg := fmt.Sprintf("%s was shot!", player)
			game.MessageChannel(msg)
			return
		}
	}
}

func (r Roulette) EndGame(id rb.Key) {
	rb.DeleteGame(id)
}

func (r Roulette) String() string {
	return "Roulette"
}

func main() {
	// no := rb.Null_Operator{}
	// rb.StartBot(no)

	rr := Roulette{}
	rb.StartBot(rr)
}
