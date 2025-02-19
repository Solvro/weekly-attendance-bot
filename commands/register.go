package commands

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func RegisterAllCommands(s *discordgo.Session) {
	if err := Register(s, BeginWeeklyCommand); err != nil {
		log.Fatalf("Unable to register /%s command: %v", BeginWeeklyCommand.Name, err)
	}

	if err := Register(s, EndWeeklyCommand); err != nil {
		log.Fatalf("Unable to register /%s command: %v", EndWeeklyCommand.Name, err)
	}
}
