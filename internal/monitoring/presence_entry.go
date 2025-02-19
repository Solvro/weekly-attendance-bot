package monitoring

import "time"

type PresenceEventType string

const PresenceJoined PresenceEventType = "joined"
const PresenceLeft PresenceEventType = "left"

type PresenceEntry struct {
	Event PresenceEventType
	At    time.Time
}
