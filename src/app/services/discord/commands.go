package discord

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"ichor-stats/src/app/models/faceit"
	"ichor-stats/src/app/services/config"
	"ichor-stats/src/app/services/discord/helpers"
	"ichor-stats/src/package/api"
	client "ichor-stats/src/package/http"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!stats") {
		requesterID := GetRequesterID(m.Author.ID)

		// Create a Bearer string by appending string access token
		var bearer = "Bearer " + config.GetConfig().FACEIT_API_KEY
		req, err := http.NewRequest("GET", api.GetFaceitPlayerCsgoStats(requesterID), nil)
		if err != nil {
			log.Println(err)
			return
		}

		// add authorization header to the req
		req.Header.Add("Authorization", bearer)

		response, err := client.Fire(req)
		if err != nil {
			log.Println(err)
			return
		}

		var stats faceit.Stats
		err = json.NewDecoder(response.Body).Decode(&stats)
		if err != nil {
			log.Println(err)
			return
		}

		req, err = http.NewRequest("GET", api.GetFaceitPlayerStats(requesterID), nil)
		if err != nil {
			log.Println(err)
			return
		}
		req.Header.Add("Authorization", bearer)

		response, err = client.Fire(req)
		if err != nil {
			log.Println(err)
			return
		}

		var user faceit.User
		err = json.NewDecoder(response.Body).Decode(&user)
		if err != nil {
			log.Println(err)
			return
		}

		embed := helpers.NewEmbed().
			SetTitle(user.Games.CSGO.Name).
			AddField("ELO", strconv.Itoa(user.Games.CSGO.ELO), true).
			AddField("Skill Level", strconv.Itoa(user.Games.CSGO.SkillLevel), true).
			AddField("Average K/D Ratio", stats.Lifetime.AverageKD, false).
			AddField("Average Headshots %", stats.Lifetime.AverageHeadshots, true)

		_, err = s.ChannelMessageSendEmbed(config.GetConfig().CHANNEL_ID, embed.MessageEmbed)
		if err != nil {
			log.Println(err)
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
