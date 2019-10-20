package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"ichor-stats/src/app/models/config"
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

	err = dc.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	} else {
		fmt.Println("Discord websocket open")
	}

	// Register discord handlers
	dc.AddHandler(MessageCreate)

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
