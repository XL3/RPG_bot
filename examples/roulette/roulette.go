package main

import (
	"fmt"
	"math/rand"

	rb "github.com/XL3/RPG_bot"
)

type Roulette struct{}

/**
 * Send a welcome message
 */
func (r Roulette) InitGame(game *rb.Game) {
	game.MessageChannel("Welcome to Roulette! Please register to start playing!")
}

/**
 * Choose a player at random to shoot
 * Mark the player that has been shot
 * Send each player what happened in a private message
 */
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

/**
 * Reveal and remove the marked player from the game's list of players
 */
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

// Send a farewell message
func (r Roulette) EndGame(game *rb.Game) {
	game.MessageChannel("Thanks for playing! Until next time!")
}

// A string representation of the game
func (r Roulette) String() string {
	return "Roulette"
}

func main() {
	// no := rb.Null_Operator{}
	// rb.StartBot(no)

	rr := Roulette{}
	rb.StartBot(rr)
}
