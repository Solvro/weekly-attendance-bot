package commands

import (
	"fmt"
	"github.com/Solvro/weekly-attendance-bot/internal/monitoring"
	"github.com/Solvro/weekly-attendance-bot/internal/results"
	"github.com/bwmarrin/discordgo"
)

var EndWeeklyCommand = &discordgo.ApplicationCommand{
	Name:                     "end-weekly",
	Description:              "End the weekly session",
	DefaultMemberPermissions: &administratorPermissions,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionChannel,
			Name:        "voice_channel",
			Description: "The voice channel to stop checking for attendance",
			Required:    true,
			ChannelTypes: []discordgo.ChannelType{
				discordgo.ChannelTypeGuildVoice,
			},
		},
	},
}

func HandleEndWeekly(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := i.ApplicationCommandData().Options
	var channel *discordgo.Channel

	for _, option := range options {
		switch option.Name {
		case "voice_channel":
			channel = option.ChannelValue(s)
		}
	}

	if channel == nil {
		return respond(s, i, "Bad channel provided")
	}

	if !monitoring.IsRecordingChannel(channel) {
		return respond(s, i, "Selected channel is not being monitored, nothing to end")
	}

	entries, minimalDuration, err := monitoring.End(channel.ID)
	if err != nil {
		return err
	}

	attendancePerUser := results.CompileFromEvents(s, channel.ID, entries, minimalDuration)

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%d users attended:\n%s", len(attendancePerUser), attendancePerUser.String()),
		},
	})
}
