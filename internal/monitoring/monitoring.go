package monitoring

import (
	"errors"
	"github.com/Solvro/weekly-attendance-bot/dtos"
	"github.com/bwmarrin/discordgo"
	"sync"
	"time"
)

type recorder struct {
	s        *discordgo.Session
	mu       sync.RWMutex
	channels map[string]channelRecorder
}

type channelRecorder struct {
	mu              sync.RWMutex
	minimalDuration time.Duration
	users           map[string][]dtos.PresenceEntry
}

var recorded = recorder{
	mu:       sync.RWMutex{},
	channels: map[string]channelRecorder{},
}

func Begin(ch *discordgo.Channel, duration uint64) {
	recorded.mu.Lock()
	defer recorded.mu.Unlock()

	recorded.channels[ch.ID] = channelRecorder{
		mu:              sync.RWMutex{},
		minimalDuration: time.Duration(duration) * time.Minute,
		users:           map[string][]dtos.PresenceEntry{},
	}
}

func IsRecordingChannel(ch *discordgo.Channel) bool {
	recorded.mu.RLock()
	defer recorded.mu.RUnlock()
	_, ok := recorded.channels[ch.ID]
	return ok
}

func End(channelID string) (map[string][]dtos.PresenceEntry, time.Duration, error) {
	recorded.mu.Lock()
	defer recorded.mu.Unlock()

	perChannel, ok := recorded.channels[channelID]
	if !ok {
		return nil, 0, errors.New("no recording was started for this channel")
	}

	defer delete(recorded.channels, channelID)

	return perChannel.users, perChannel.minimalDuration, nil
}
