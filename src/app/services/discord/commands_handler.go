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

		var user faceit.User
		_ = GetFaceitStats(api.GetFaceitPlayerStats(requesterID)).Decode(&user)

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

func HandleParameterisedCommand(requesterId string, command []string, user faceit.User, messages *[]*helpers.Embed) {
	switch trimLeftChar(command[0]) {
		case "map": EmbeddedMapStats(requesterId, user, command[1], messages)
		case "last": EmbeddedLastMatchStats(requesterId, user, command[1], messages)
		case "totals": EmbeddedLastMatchTotals(requesterId, user, command[1], messages)
	}
}

func HandleCommand(requesterId string, command string, user faceit.User, messages *[]*helpers.Embed) {
	switch trimLeftChar(command) {
		case "stats": EmbeddedStats(requesterId, user, messages)
		case "streak": EmbeddedStreak(requesterId, user, messages)
		case "recent": EmbeddedStreak(requesterId, user, messages)
		case "green": EmbeddedGreen(messages)
	}
}

func GetRequesterID(discordID string) string {
	if discordID == "210457267710066689" {
		return "0d94613d-b736-46ba-b8cd-d2159ddad705" // Tecey
	} else if discordID == "210449893892947969" {
		return "b26df7d4-8517-4ec6-ab58-708487e5fe60" // Dylan
	} else if discordID == "210438278623526913" {
		return "b0a57a5a-2f7a-481c-aaa8-8013a83378e3" // Jamie
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