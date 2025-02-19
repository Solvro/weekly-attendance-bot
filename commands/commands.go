package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

var administratorPermissions int64 = discordgo.PermissionAdministrator

func Register(s *discordgo.Session, cmd *discordgo.ApplicationCommand) error {
	_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
	fmt.Println("here")
	return err
}

func Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var handler func(s *discordgo.Session, i *discordgo.InteractionCreate) error

	switch i.ApplicationCommandData().Name {
	case "begin-weekly":
		handler = HandleBeginWeekly
	case "end-weekly":
		handler = HandleEndWeekly
	default:
		_ = respond(s, i, "Invalid slash command")
		return
	}

	if err := handler(s, i); err != nil {
		log.Printf("Couldn't respond to slash command: %v", err)
	}
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, message string) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}
