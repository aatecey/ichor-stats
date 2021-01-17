package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mitchellh/mapstructure"
	configModel "ichor-stats/src/app/models/config"
	"ichor-stats/src/app/models/players"
	"ichor-stats/src/app/services/config"
	"ichor-stats/src/package/discord"
	"log"
	"strings"
)

func NewDiscordHandler(ds *ServiceDiscord, config *configModel.Configuration) {
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

	dc.AddHandler(MessageCreate)
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!") {
		playerName := GetRequesterID(m.Author.ID)

		command := strings.TrimSpace(m.Content)
		commandString := strings.Split(command, " ")

		var messages = make([]*discord.Embed, 0)

		if len(commandString) == 1 {
			HandleCommand(playerName, commandString[0], &messages)
		} else {
			HandleParameterisedCommand(playerName, commandString, &messages)
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
			return playerDetails.Name
		}
	}

	return ""
}