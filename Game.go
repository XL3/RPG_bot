package RPG_bot

import (
	"errors"
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
)

type Game struct {
	ID      Key
	ChID    discord.ChannelID
	Players []*Player
}

var active_games map[Key]*Game

func (game *Game) RegisterPlayer(userID string) *Player {
	for _, p := range game.Players {
		if p.id.String() == userID {
			return p
		}
	}

	player := createPlayer(userID)
	game.Players = append(game.Players, &player)
	return &player
}

func (game *Game) String() string {
	var str string
	if game.ID > 0 {
		str = fmt.Sprintf("`Game[%d]:` ", game.ID)
	}

	return str
}

func (game *Game) MessagePlayer(player Player, message string) error {
	id := discord.UserID(player.id)
	dm, err := bot_session.CreatePrivateChannel(id)

	if err == nil {
		_, err = bot_session.SendMessage(dm.ID, game.String()+message, nil)
	}
	if err != nil {
		log.Fatal("Failed to message player ", err)
	}
	return err
}

func (game *Game) MessageChannel(message string) error {
	_, err := bot_session.SendMessage(game.ChID, game.String()+message, nil)
	return err
}

func (game *Game) respondToInteraction(e *gateway.InteractionCreateEvent, response string) error {
	// Respond to interaction
	data := api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content: game.String() + response,
		},
	}

	err := bot_session.RespondInteraction(e.ID, e.Token, data)
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

func DeleteGame(id Key) (*Game, error) {
	game, err := GetGame(id)
	if err == nil {
		delete(active_games, id)
	}

	return game, err
}
