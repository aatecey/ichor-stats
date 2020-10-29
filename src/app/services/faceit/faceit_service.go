package faceit

import (
	"encoding/json"
	"ichor-stats/src/app/models/config"
	"ichor-stats/src/app/models/faceit"
	"ichor-stats/src/app/models/players"
	"ichor-stats/src/app/services/discord"
	"ichor-stats/src/app/services/discord/helpers"
	"ichor-stats/src/package/api"
	client "ichor-stats/src/package/http"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type ServiceFaceit struct {
	Config config.Configuration
	DiscordService discord.ServiceDiscord
}

func NewFaceitService(config *config.Configuration, ds discord.ServiceDiscord) ServiceFaceit {
	return ServiceFaceit{
		Config: *config,
		DiscordService: ds,
	}
}

func (fs *ServiceFaceit) MatchEnd(playerFromMatch faceit.Players, messages *[]*helpers.Embed, stats faceit.Match) {
	var outcome = "Victory"

	if playerFromMatch.Stats.Result == "0" {
		outcome = "Defeat"
	}

	*messages = append(*messages, helpers.NewEmbed().
		SetTitle("Match ended for "+playerFromMatch.Nickname).
		SetDescription(outcome+" on "+stats.Rounds[0].MatchStats.Map+" ["+
			stats.Rounds[0].MatchStats.Score+"]").
		AddField("Kills", playerFromMatch.Stats.Kills, true).
		AddField("Assists", playerFromMatch.Stats.Assists, true).
		AddField("Deaths", playerFromMatch.Stats.Deaths, true).
		AddField("K/D Ratio", playerFromMatch.Stats.KD, true).
		AddField("K/R Ratio", playerFromMatch.Stats.KR, true))
}

func (fs *ServiceFaceit) MatchCreated(webhook faceit.Webhook, messages *[]*helpers.Embed) {
	req, err := http.NewRequest("GET", api.GetFaceitMatch(webhook.Payload.MatchID), nil)
	req.Header.Add("Authorization", "Bearer " + fs.Config.FACEIT_API_KEY)
	response, err := client.Fire(req)
	body, err := ioutil.ReadAll(response.Body)

	log.Println("Match Created")
	log.Println(string(body))

	var stats faceit.Match
	_ = json.Unmarshal(body, &stats)

	for _, s := range stats.Rounds {
		for _, a := range s.Teams {
			for _, d := range a.Players {
				if _, playerPresentInMap := players.Players[d.ID]; playerPresentInMap {
					*messages = append(*messages, helpers.NewEmbed().
						SetTitle("Match created for "+d.Nickname))
				}
			}
		}
	}

	if err != nil {
		log.Println(err)
	}
}

func (fs *ServiceFaceit) MatchReady(webhook faceit.Webhook, messages *[]*helpers.Embed) {
	log.Println("Match Ready")

	var message = helpers.NewEmbed()

	for _, team := range webhook.Payload.MatchTeams {
		var messageValue = ""

		for _, player := range team.Roster {
			messageValue = messageValue + "Level " + strconv.Itoa(player.SkillLevel) + "\t- " + player.Nickname + "\n"

			if _, playerPresentInMap := players.Players[player.ID]; playerPresentInMap {
				message.SetTitle("Match Created for " + player.Nickname)
			}
		}

		message.AddField("Team " + team.Name[5:len(team.Name)], messageValue, false)
	}

	*messages = append(*messages, message)
}

func (fs *ServiceFaceit) MatchConfiguring(webhook faceit.Webhook, messages *[]*helpers.Embed) {
	req, err := http.NewRequest("GET", api.GetFaceitMatch(webhook.Payload.MatchID), nil)
	req.Header.Add("Authorization", "Bearer " + fs.Config.FACEIT_API_KEY)
	response, err := client.Fire(req)
	body, err := ioutil.ReadAll(response.Body)

	log.Println("Match Configuring")
	log.Println(string(body))

	var stats faceit.Match
	_ = json.Unmarshal(body, &stats)

	for _, s := range stats.Rounds {
		for _, a := range s.Teams {
			for _, d := range a.Players {
				if _, playerPresentInMap := players.Players[d.ID]; playerPresentInMap {
					*messages = append(*messages, helpers.NewEmbed().
						SetTitle("Match configuring for "+d.Nickname))
				}
			}
		}
	}

	if err != nil {
		log.Println(err)
	}
}