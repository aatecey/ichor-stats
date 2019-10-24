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
	"strings"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!") {
		requesterID := GetRequesterID(m.Author.ID)

		var stats faceit.Stats
		var user faceit.User

		statsErr := GetFaceitStats(api.GetFaceitPlayerCsgoStats(requesterID)).Decode(&stats)
		if statsErr != nil {
			log.Println(statsErr)
		}

		userErr := GetFaceitStats(api.GetFaceitPlayerStats(requesterID)).Decode(&user)
		if userErr != nil {
			log.Println(userErr)
		}

		command := strings.TrimSpace(m.Content)
		commandString := strings.Split(command, " ")

		var embed *helpers.Embed
		if len(commandString) == 1 {
			embed = HandleCommand(commandString[0], stats, user, requesterID)
		} else {
			embed = HandleParameterisedCommand(commandString, stats, user, requesterID, s)
		}

		if embed != nil {
			_, err := s.ChannelMessageSendEmbed(config.GetConfig().CHANNEL_ID, embed.MessageEmbed)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func GetFaceitStats(apiUrl string) *json.Decoder  {
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + config.GetConfig().FACEIT_API_KEY
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Println(err)
		return nil
	}

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	response, err := client.Fire(req)
	if err != nil {
		log.Println(err)
		return nil
	}

	return json.NewDecoder(response.Body)
}

func HandleParameterisedCommand(command []string, stats faceit.Stats, user faceit.User, requesterID string, session *discordgo.Session) *helpers.Embed {
	switch trimLeftChar(command[0]) {
		case "map": return EmbeddedMapStats(stats, user, command[1])
		case "last": return EmbeddedLastFiveStats(stats, user, command[1], requesterID, session)
	}

	return nil
}

func HandleCommand(command string, stats faceit.Stats, user faceit.User, requesterId string) *helpers.Embed {
	switch trimLeftChar(command) {
		case "stats": return EmbeddedStats(stats, user)
		case "streak": return EmbeddedStreak(stats, user)
		case "recent": return EmbeddedStreak(stats, user)
		case "green": return EmbeddedGreen()
	}

	return nil
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

func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}