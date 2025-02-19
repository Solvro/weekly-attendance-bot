package monitoring

import (
	"github.com/Solvro/weekly-attendance-bot/dtos"
	"github.com/Solvro/weekly-attendance-bot/internal/config"
	"github.com/Solvro/weekly-attendance-bot/internal/storage"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

func ForceJoin(channelID string, userID string) {
	recorded.mu.RLock()
	defer recorded.mu.RUnlock()

	joined, ok := recorded.channels[channelID]
	if !ok {
		return
	}

	joined.mu.Lock()
	defer joined.mu.Unlock()

	at := time.Now()
	joined.users[userID] = append(joined.users[userID], dtos.PresenceEntry{
		Event: dtos.PresenceJoined,
		At:    at,
	})

	if err := storage.InsertEvent(channelID, userID, dtos.PresenceJoined, at); err != nil {
		log.Printf("Failed inserting an event into the db: %v", err)
	}

	if config.Logging {
		log.Printf("Force joined user %s to be monitored in channel %s\n", userID, channelID)
	}
}

func ProcessEvents(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	recorded.mu.RLock()
	defer recorded.mu.RUnlock()

	at := time.Now()

	joined, isMonitoringJoined := recorded.channels[vs.ChannelID]
	if vs.BeforeUpdate == nil {
		if !isMonitoringJoined {
			return
		}

		joined.mu.Lock()
		defer joined.mu.Unlock()

		joined.users[vs.UserID] = append(joined.users[vs.UserID], dtos.PresenceEntry{
			Event: dtos.PresenceJoined,
			At:    at,
		})

		if err := storage.InsertEvent(vs.ChannelID, vs.UserID, dtos.PresenceJoined, at); err != nil {
			log.Printf("Failed inserting an event into the db: %v", err)
		}

		if config.Logging {
			log.Printf("User %v joined a monitored channel (%v)\n", vs.UserID, vs.ChannelID)
		}

		return
	}

	left, isMonitoringLeft := recorded.channels[vs.BeforeUpdate.ChannelID]

	if !isMonitoringLeft && !isMonitoringJoined {
		// channel not being recorded at all
		return
	}

	// user left a monitored channel
	if isMonitoringLeft && vs.BeforeUpdate.ChannelID != vs.ChannelID {
		left.mu.Lock()
		defer left.mu.Unlock()
		left.users[vs.UserID] = append(left.users[vs.UserID], dtos.PresenceEntry{
			Event: dtos.PresenceLeft,
			At:    at,
		})

		if err := storage.InsertEvent(vs.ChannelID, vs.UserID, dtos.PresenceLeft, at); err != nil {
			log.Printf("Failed inserting an event into the db: %v", err)
		}

		if config.Logging {
			log.Printf("User %v left a monitored channel (%v)\n", vs.UserID, vs.BeforeUpdate.ChannelID)
		}
	}

	// user joined a monitored channel from some other voice channel
	if isMonitoringJoined && vs.BeforeUpdate.ChannelID != vs.ChannelID {
		joined.mu.Lock()
		defer joined.mu.Unlock()

		joined.users[vs.UserID] = append(joined.users[vs.UserID], dtos.PresenceEntry{
			Event: dtos.PresenceJoined,
			At:    at,
		})

		if err := storage.InsertEvent(vs.ChannelID, vs.UserID, dtos.PresenceJoined, at); err != nil {
			log.Printf("Failed inserting an event into the db: %v", err)
		}

		if config.Logging {
			log.Printf("User %v joined a monitored channel (%v) by switching from (%v)\n", vs.UserID, vs.ChannelID, vs.BeforeUpdate.ChannelID)
		}
	}
}
