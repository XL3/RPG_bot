package RPG_bot

import (
	"log"
	"os"
	"strings"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/session"
	"github.com/diamondburned/arikawa/v2/utils/httputil/httpdriver"
	_ "github.com/joho/godotenv/autoload"
)

var bot_session *session.Session

func createSession() (discord.AppID, discord.GuildID) {
	mustSnowflakeEnv := func(env string) discord.Snowflake {
		s, err := discord.ParseSnowflake(os.Getenv(env))
		if err != nil {
			log.Fatalf("Invalid snowflake for $%s: %v", env, err)
		}
		return s
	}

	appID := discord.AppID(mustSnowflakeEnv("APP_ID"))
	guildID := discord.GuildID(mustSnowflakeEnv("GUILD_ID"))

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatalln("No $BOT_TOKEN given.")
	}

	s, err := session.New("Bot " + token)
	if err != nil {
		log.Fatalln("Session failed:", err)
	}

	bot_session = s
	return appID, guildID
}

func configureSession(opr Operator) {
	configureCommandHandlers(opr)

	bot_session.AddHandler(func(e *gateway.InteractionCreateEvent) {
		command := e.Data.Name
		if handler, ok := rpg_commands[command]; ok {
			handler.process(e)
		}
	})

	bot_session.Gateway.AddIntents(gateway.IntentGuilds)
	bot_session.Gateway.AddIntents(gateway.IntentGuildMessages)
	bot_session.Gateway.AddIntents(gateway.IntentDirectMessages)
	bot_session.Gateway.AddIntents(gateway.IntentGuildMembers)
}

func createGuildCommands(appID discord.AppID, guildID discord.GuildID) {
	commands, err := bot_session.GuildCommands(appID, guildID)
	if err != nil {
		log.Fatalln("failed to get guild commands:", err)
	}

	for _, command := range commands {
		log.Println("Existing command", command.Name, command.ID, "found.")
	}

	newCommands := []api.CreateCommandData{}
	for _, command := range rpg_commands {
		newCommands = append(newCommands, command.command)
	}
	// newCommands = append(newCommands, bonk_command)

	for _, command := range newCommands {
		_, err := bot_session.CreateGuildCommand(appID, guildID, command)

		if err != nil {
			log.Fatalln("failed to create guild command:", err)
		}
	}
}

func getGuildMembersAbbrev(appID discord.AppID, guildID discord.GuildID) (am []Member) {
	url := api.EndpointGuilds + guildID.String() + "/members"

	limit_opt := func(r httpdriver.Request) error {
		r.AddQuery(map[string][]string{
			"limit": {"10"},
		})
		return nil
	}

	var members []discord.Member
	if err := bot_session.RequestJSON(&members, "GET", url, limit_opt); err != nil {
		log.Fatalln("failed to get guild members:", err)
	}

	for _, mem := range members {
		am = append(am, Member{
			id:   discord.Snowflake(mem.User.ID),
			name: mem.User.Username,
			nick: mem.Nick,
		})
	}
	return
}

func getUserID(e *gateway.InteractionCreateEvent, name string) string {
	userID := e.Member.User.ID.String()

	if strings.HasPrefix(name, "<@!") {
		userID = name[3 : len(name)-1]
	} else if len(name) > 0 {
		for _, mem := range guild_members {
			if name == mem.name || name == mem.nick {
				userID = mem.id.String()
				break
			}
		}
	}

	return userID
}

type Member struct {
	id   discord.Snowflake
	name string
	nick string
}

var guild_members []Member

func StartBot(opr Operator) {
	appID, guildID := createSession()
	configureSession(opr)

	if err := bot_session.Open(); err != nil {
		log.Fatalln("failed to open:", err)
	}
	defer bot_session.Close()
	log.Println("Gateway connected.")

	createGuildCommands(appID, guildID)
	guild_members = getGuildMembersAbbrev(appID, guildID)

	// Block forever.
	select {}
}
