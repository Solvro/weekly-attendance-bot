## A Discord bot for tracking user presence

### Usage
```bash
go build -o bot cmd/bot/main.go
./bot [options] 
```

Options:  
- `-token` (required) - bot's secret  
- `-enable-logging` - whether to display non-essential logging

### Commands
- `/begin-weekly <channel> <duration=5>`
  - channel - voice channel in which the users' presence will be recorded
  - duration - minimal timespan for which the user must be present in the given channel for their attendance to be recorded  (in seconds)
- `end-weekly <channel>`

### Notes
Due to how `discordgo` (and probably the whole API) works, it's required that the bot is started before users join
the channel that is to be monitored for their presence to be tracked.  
I have not found a way to retrieve a set of users in a given voice channel without the join events being 
received - that's why the bot must be running at the point they are sent.


#### **The implementation has not been stress tested!**
