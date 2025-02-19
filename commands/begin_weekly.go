package commands

import (
	"fmt"
	"github.com/Solvro/weekly-attendance-bot/internal/monitoring"
	"github.com/bwmarrin/discordgo"
)

var BeginWeeklyCommand = &discordgo.ApplicationCommand{
	Name:                     "begin-weekly",
	Description:              "Start a weekly meeting",
	DefaultMemberPermissions: &administratorPermissions,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionChannel,
			Name:        "voice_channel",
			Description: "The voice channel to check for attendance",
			Required:    true,
			ChannelTypes: []discordgo.ChannelType{
				discordgo.ChannelTypeGuildVoice,
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "duration",
			Description: "Time period user needs to spend in the meeting for presence to be marked (in minutes, default: 5)",
			Required:    false,
		},
	},
}

func HandleBeginWeekly(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := i.ApplicationCommandData().Options
	var channel *discordgo.Channel
	var durationToRecord uint64 = 5

	for _, option := range options {
		switch option.Name {
		case "voice_channel":
			channel = option.ChannelValue(s)
		case "duration":
			durationToRecord = option.UintValue()
		}
	}

	if channel == nil {
		return respond(s, i, "Bad channel provided")
	}

	if monitoring.IsRecordingChannel(channel) {
		return respond(s, i, "Channel already being recorded, end the monitoring first by using the `/end-weekly` command")
	}

	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		return err
	}

	monitoring.Begin(channel, durationToRecord)

	// record initial users (no join event will be transmitted for them)
	// WARNING: bot must start BEFORE first users join the channel,
	// otherwise their presence won't be kept in the state, and they will have to rejoin
	for _, vs := range guild.VoiceStates {
		if vs.ChannelID == channel.ID {
			monitoring.ForceJoin(channel.ID, vs.UserID)
		}
	}

	return respond(s, i, fmt.Sprintf("Attendance monitoring started for channel %s, minimal attendance: %d minutes", channel.Name, durationToRecord))
}
