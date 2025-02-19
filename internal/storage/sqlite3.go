package storage

import (
	"database/sql"
	"github.com/Solvro/weekly-attendance-bot/dtos"
	"github.com/Solvro/weekly-attendance-bot/internal/config"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

var DB *sql.DB

func ConnectToSqlite() error {
	var err error
	DB, err = sql.Open("sqlite3", config.DatabasePath)
	return err
}

const CreateEventsTableQuery = `
create table if not exists events (
    id integer primary key autoincrement,
    channel_id text not null,
    user_id text not null,
    kind text not null,
    at text not null
)`

const InsertEventQuery = `
insert into events(channel_id, user_id, kind, at) values(?, ?, ?, ?)
`

func RunInitialMigrations() error {
	_, err := DB.Exec(CreateEventsTableQuery)
	return err
}

func InsertEvent(channelID string, userID string, kind dtos.PresenceEventType, at time.Time) error {
	_, err := DB.Exec(InsertEventQuery, channelID, userID, kind, at)
	if err == nil && config.Logging {
		log.Printf("inserted event (%s, %s, %s, %s) in the underlying storage\n", channelID, userID, kind, at)
	}
	return err
}
