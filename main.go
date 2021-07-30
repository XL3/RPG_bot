package main

import (
	"log"
)

func main() {
	appID, guildID := createSession()
	configureSession()

	if err := bot_session.Open(); err != nil {
		log.Fatalln("failed to open:", err)
	}
	defer bot_session.Close()
	log.Println("Gateway connected.")

	createGuildCommands(appID, guildID)
	guild_members = abbrevMembers(getGuildMembers(appID, guildID))
	log.Println(guild_members)

	// Block forever.
	select {}
}
