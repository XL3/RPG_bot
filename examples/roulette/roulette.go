package main

import (
	"fmt"
	"log"
	"math/rand"
	rb "rpg_bot"
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
			msg = fmt.Sprintf("game[%d] Spared", game.ID)
		} else {
			player.Role = "shot"
			msg = fmt.Sprintf("game[%d] SHOT", game.ID)
			player.Score = -1
		}

		if err := game.MessagePlayer(player, msg); err != nil {
			log.Fatal("Failed to message player", err)
		}
	}
}

func (r Roulette) UpdateState(game *rb.Game) {
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
	rr := Roulette{}
	log.Println(rr)

	rb.StartBot(rr)
}
