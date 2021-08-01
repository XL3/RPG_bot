package rpg_bot

import (
	"errors"
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
)

type Game struct {
	ID      Key
	ChID    discord.ChannelID
	Players []Player
}

var active_games map[Key]*Game

func (game *Game) RegisterPlayer(userID string) {
	for _, p := range game.Players {
		if p.id.String() == userID {
			return
		}
	}

	player := createPlayer(userID)
	game.Players = append(game.Players, player)
}

func (game *Game) PrefixMessage(message string) string {
	return fmt.Sprintf("`Game[%d]:` %s", game.ID, message)
}

func (game *Game) MessagePlayer(player Player, message string) error {
	id := discord.UserID(player.id)
	dm, err := bot_session.CreatePrivateChannel(id)

	if err == nil {
		_, err = bot_session.SendMessage(dm.ID, game.PrefixMessage(message), nil)
	}
	return err
}

func (game *Game) MessageChannel(message string) error {
	_, err := bot_session.SendMessage(game.ChID, game.PrefixMessage(message), nil)
	return err
}

func GetGame(id Key) (*Game, error) {
	if active_games == nil {
		active_games = make(map[Key]*Game)
	}
	game, ok := active_games[id]
	if !ok {
		return nil, errors.New("game not found")
	}
	return game, nil
}

func CreateNewGame(id Key) (*Game, error) {
	game, err := GetGame(id)
	if err == nil {
		return game, errors.New("game id exists")
	}

	active_games[id] = &Game{ID: id}
	return active_games[id], nil
}

func DeleteGame(id Key) {
	_, err := GetGame(id)

	if err == nil {
		delete(active_games, id)
	}
}
