# RPG_bot
A discord bot in Go to manage generic role-playing games.

## Why?
I wanted to practice Go, and make something fun that people can tinker with.

## How does it work?
Go to the [Discord Developer Portal](https://discord.com/developers/) and get your bot set up.
The bot's environment should contain the variables `BOT_TOKEN`, `APP_ID`, and `GUILD_ID`. 

The idea is to implement the `Operator` interface
```go
type Operator interface {
	InitGame(game *Game)
	StartTurn(game *Game)
	EndTurn(game *Game)
	EndGame(game *Game)
	String() string
}
```

Think of these functions as callbacks the bot calls when it receives a command. These callbacks
constitute the behavior of your game. The `Operator` controls the game, sort of like a dungeon master.

Perhaps you would like to send a welcome message to players when you start the game, or perhaps you
would like to assign them roles in private messages. How you build your game is entirely up to you,
and an example is provided in [roulette](examples/roulette/roulette.go).

A game would start by initializing it in a channel and giving it its `id`. Each player would
register themselves to the game. Afterwards turns are started and ended until a decision is made to
end the game. All through slash commands.