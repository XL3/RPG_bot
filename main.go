package main

import (
	"log"
)

func main() {
	appID, guildID := create_session()
	configure_session()

	if err := bot_session.Open(); err != nil {
		log.Fatalln("failed to open:", err)
	}
	defer bot_session.Close()
	log.Println("Gateway connected.")

	create_guild_commands(appID, guildID)
	guild_members = abbrev_members(get_guild_members(appID, guildID))

	// Block forever.
	select {}
}
