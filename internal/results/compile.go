package results

import (
	"bytes"
	"fmt"
	"github.com/Solvro/weekly-attendance-bot/dtos"
	"github.com/Solvro/weekly-attendance-bot/internal/storage"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

type Compiled map[*discordgo.User]time.Duration

func CompileFromEvents(s *discordgo.Session, channelID string, entriesPerUser map[string][]dtos.PresenceEntry, minimalDuration time.Duration) Compiled {
	attendancePerUser := map[*discordgo.User]time.Duration{}
	at := time.Now()
	for userId, entries := range entriesPerUser {
		user, err := s.User(userId)
		if err != nil {
			log.Printf("There was an issue with fetching details of user %v, this user won't be included in the final report", userId)
			continue
		}

		if err := storage.InsertEvent(channelID, userId, dtos.PresenceLeft, at); err != nil {
			log.Printf("Failed inserting an event into the db: %v", err)
		}

		currentStatus := dtos.PresenceLeft
		var at time.Time
		totalDuration := 0 * time.Minute
		for _, entry := range entries {
			// this shouldn't happen at all but if it occurs, we are treating this case
			// for the benefit of the user, e.g. user joins at 10:00, but then also joins at 10:10
			// (maybe their session got corrupted, bot hang up, etc.) it is treated as if the user
			// only joined at 10:00 and never exited the channel afterward
			if entry.Event == currentStatus {
				continue
			}

			if entry.Event == dtos.PresenceJoined {
				at = entry.At
			} else if currentStatus == dtos.PresenceJoined && entry.Event == dtos.PresenceLeft {
				totalDuration += entry.At.Sub(at)
			}
			currentStatus = entry.Event
		}

		if currentStatus == dtos.PresenceJoined {
			totalDuration += time.Now().Sub(at)
		}

		if totalDuration >= minimalDuration {
			attendancePerUser[user] = totalDuration
		}
	}

	return attendancePerUser
}

func (c Compiled) String() string {
	var buffer bytes.Buffer
	for user, duration := range c {
		buffer.WriteString(fmt.Sprintf("- %s (%s): %s", user.Username, user.GlobalName, duration.Round(time.Second)))
	}
	return buffer.String()
}
