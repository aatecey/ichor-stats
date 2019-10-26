package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	configModel "ichor-stats/src/app/models/config"
	"ichor-stats/src/app/models/faceit"
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
	dc.AddHandler(dh.MessageCreate)

	return dh
}

func (dh *HandlerDiscord) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!") {
		requesterID := GetRequesterID(m.Author.ID)

		var stats faceit.Stats
		var user faceit.User

		statsErr := api.FaceitRequest(api.GetFaceitPlayerCsgoStats(requesterID)).Decode(&stats)
		if statsErr != nil {
			log.Println(statsErr)
			return
		}

		userErr := api.FaceitRequest(api.GetFaceitPlayerStats(requesterID)).Decode(&user)
		if userErr != nil {
			log.Println(userErr)
			return
		}

		command := strings.TrimSpace(m.Content)
		commandString := strings.Split(command, " ")

		var embed *helpers.Embed
		if len(commandString) == 1 {
			embed = dh.ServiceDiscord.HandleCommand(commandString[0], stats, user, requesterID)
		} else {
			embed = dh.ServiceDiscord.HandleParameterisedCommand(commandString, stats, user, requesterID, s)
		}

		if embed != nil {
			_, err := s.ChannelMessageSendEmbed(config.GetConfig().CHANNEL_ID, embed.MessageEmbed)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func GetRequesterID(discordID string) string {
	if discordID == "210457267710066689" {
		return "0d94613d-b736-46ba-b8cd-d2159ddad705"
	} else if discordID == "210449893892947969" {
		return "b26df7d4-8517-4ec6-ab58-708487e5fe60"
	} else if discordID == "210438278623526913" {
		return "b0a57a5a-2f7a-481c-aaa8-8013a83378e3"
	}

	return ""
}
