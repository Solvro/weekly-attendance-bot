package main

import (
	"github.com/Solvro/weekly-attendance-bot/commands"
	"github.com/Solvro/weekly-attendance-bot/internal/config"
	"github.com/Solvro/weekly-attendance-bot/internal/monitoring"
	"github.com/bwmarrin/discordgo"
	"log"
)

func init() {
	config.LoadAndValidate()
}

func main() {
	dg, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	dg.StateEnabled = true
	dg.Identify.Intents = discordgo.IntentsGuildVoiceStates

	dg.AddHandler(monitoring.ProcessEvents)
	dg.AddHandler(commands.Handle)

	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening Discord session: %v", err)
	}

	commands.RegisterAllCommands(dg)

	log.Printf("Bot is running")
	select {}
}
