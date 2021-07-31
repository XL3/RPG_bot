package main

import (
	"fmt"
	"log"
	"math/rand"
	rb "rpg_bot"
)

type Roulette struct{}

func (r Roulette) InitGame(id rb.Key) *rb.Game {
	session, err := rb.CreateNewGame(id)

	if err != nil {
		log.Fatal("Failed to initialize game", err)
	}

	return session
}

func (r Roulette) StartTurn(game *rb.Game) {
	size := len(game.Players)
	target := rand.Intn(size)

	for i, player := range game.Players {
		if i != target {
			role := fmt.Sprintf("game[%d] Spared", game.ID)
			player.AssignAndDMRole(role)
		} else {
			role := fmt.Sprintf("game[%d] Shot", game.ID)
			player.AssignAndDMRole(role)
			player.Score = -1
		}

		log.Println(game.Players[i])
	}
}

func (r Roulette) UpdateState(game *rb.Game) {
	size := len(game.Players)
	for i, player := range game.Players {
		if player.Score < 0 {
			game.Players[i] = game.Players[size-1]
			game.Players = game.Players[:size-1]
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
