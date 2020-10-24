package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mitchellh/mapstructure"
	configModel "ichor-stats/src/app/models/config"
	"ichor-stats/src/app/models/faceit"
	"ichor-stats/src/app/models/players"
	"ichor-stats/src/app/services/config"
	"ichor-stats/src/app/services/discord/helpers"
	"ichor-stats/src/package/api"
	"log"
	"strings"
)

type HandlerDiscord struct {
	Config configModel.Configuration
	ServiceDiscord ServiceDiscord
}

func NewDiscordHandler(ds *ServiceDiscord, config *configModel.Configuration) HandlerDiscord {
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

	dh := HandlerDiscord{
		Config: *config,
		ServiceDiscord: *ds,
	}
	// Register discord handlers
	dc.AddHandler(MessageCreate)

	return dh
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!") {
		requesterID := GetRequesterID(m.Author.ID)

		var user faceit.User
		_ = api.FaceitRequest(api.GetFaceitPlayerStats(requesterID)).Decode(&user)

		command := strings.TrimSpace(m.Content)
		commandString := strings.Split(command, " ")

		var messages = make([]*helpers.Embed, 0)

		if len(commandString) == 1 {
			HandleCommand(requesterID, commandString[0], user, &messages)
		} else {
			HandleParameterisedCommand(requesterID, commandString, user, &messages)
		}

		if len(messages) > 0 {
			for _, message := range messages {
				_, err := s.ChannelMessageSendEmbed(config.GetConfig().CHANNEL_ID, message.MessageEmbed)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func GetRequesterID(discordID string) string {
	for player := range players.Players {
		var playerDetails players.PlayerDetails
		_ = mapstructure.Decode(players.Players[player], &playerDetails)

		if playerDetails.DiscordId == discordID {
			return player
		}
	}

	return ""
}