package endpoints

import (
	"encoding/json"
	"github.com/labstack/echo"
	"ichor-stats/src/app/models/message"
	"ichor-stats/src/app/models/players"
	"ichor-stats/src/app/services/config"
	servicesDiscord "ichor-stats/src/app/services/discord"
	"ichor-stats/src/package/discord"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type MessageEndpointHandler struct {
	DiscordService servicesDiscord.ServiceDiscord
}

func (fh *MessageEndpointHandler) Init(echo *echo.Echo, ds servicesDiscord.ServiceDiscord) {
	singlePlayerGroup := echo.Group("/message")
	singlePlayerGroup.POST("/match-end", fh.MatchEndMessage)
	singlePlayerGroup.POST("/match-ready", fh.MatchReadyMessage)
}

func (fh *MessageEndpointHandler) MatchEndMessage(context echo.Context) error {
	body, err := ioutil.ReadAll(context.Request().Body)

	var message message.Match
	err = json.Unmarshal(body, &message)

	discordMessage := discord.NewEmbed().
		SetTitle("Match ended for " + message.Player).
		SetDescription(message.Result + " on " + message.Map + " [" + message.Score + "]").
		AddField("Kills", message.Kills, true).
		AddField("Assists", message.Assists, true).
		AddField("Deaths", message.Deaths, true).
		AddField("K/D Ratio", message.KillDeathRatio, true).
		AddField("K/R Ratio", message.KillRoundRatio, true)

	_, err = fh.DiscordService.Discord.ChannelMessageSendEmbed(config.GetConfig().CHANNEL_ID, discordMessage.MessageEmbed)
	if err != nil {
		log.Println(err)
	}

	return context.JSON(http.StatusOK, "")
}

func (fh *MessageEndpointHandler) MatchReadyMessage(context echo.Context) error {
	body, err := ioutil.ReadAll(context.Request().Body)

	var matchPlayers = ""
	var discordMessage = discord.NewEmbed()
	var webhookData message.Webhook
	err = json.Unmarshal(body, &webhookData)

	for _, team := range webhookData.Payload.MatchTeams {
		var messageValue = ""

		for _, player := range team.Roster {
			messageValue = messageValue + "Level " + strconv.Itoa(player.SkillLevel) + "\t- " + player.Nickname + "\n"

			if _, playerPresentInMap := players.Players[player.ID]; playerPresentInMap {
				if matchPlayers == "" {
					matchPlayers = player.Nickname
				} else {
					matchPlayers = matchPlayers + ", " + player.Nickname
				}
			}
		}

		discordMessage.AddField("Team "+team.Name[5:len(team.Name)], messageValue, false)
	}

	discordMessage.SetTitle("Match Created for " + matchPlayers)

	_, err = fh.DiscordService.Discord.ChannelMessageSendEmbed(config.GetConfig().CHANNEL_ID, discordMessage.MessageEmbed)
	if err != nil {
		log.Println(err)
	}

	return context.JSON(http.StatusOK, "")
}