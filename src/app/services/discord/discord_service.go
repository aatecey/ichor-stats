package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"ichor-stats/src/app/models/config"
	"ichor-stats/src/app/services/calls"
	"ichor-stats/src/package/discord"
)

type ServiceDiscord struct {
	Config config.Configuration
	Discord discordgo.Session
}

func NewDiscordService(config *config.Configuration) ServiceDiscord {
	dc, err := discordgo.New("Bot " + config.DISCORD_BOT_ID)
	if err != nil {
		fmt.Println(err)
	}

	return ServiceDiscord{
		Config: *config,
		Discord: *dc,
	}
}

func (ds *ServiceDiscord) SendMessage(message string) {
	_, err := ds.Discord.ChannelMessageSend(ds.Config.CHANNEL_ID, message)
	if err != nil {
		fmt.Println(err)
	}
}

func HandleParameterisedCommand(playerName string, command []string, messages *[]*discord.Embed) {
	switch trimLeftChar(command[0]) {
	case "map":
		calls.MapStats(playerName, command[1], messages)
	case "last":
		calls.LastMatchStats(playerName, command[1], "false", messages)
	case "totals":
		calls.LastMatchTotals(playerName, command[1], "false", messages)
	}
}

func HandleCommand(playerName string, command string, messages *[]*discord.Embed) {
	switch trimLeftChar(command) {
	case "stats":
		calls.Stats(playerName, messages)
	case "streak":
		calls.Streak(playerName, "false", messages)
	case "recent":
		calls.Streak(playerName, "false", messages)
	case "green":
		calls.Green(messages)
	}
}

func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}