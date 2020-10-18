package faceit

import (
	"encoding/json"
	"ichor-stats/src/app/models/config"
	"ichor-stats/src/app/models/faceit"
	"ichor-stats/src/app/services/discord"
	"ichor-stats/src/app/services/discord/helpers"
	"ichor-stats/src/package/api"
	client "ichor-stats/src/package/http"
	"io/ioutil"
	"log"
	"net/http"
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

func (fs *ServiceFaceit) MatchEnd(webhook faceit.Webhook, messages *[]*helpers.Embed) {
	req, err := http.NewRequest("GET", api.GetFaceitMatch(webhook.Payload.MatchID), nil)
	req.Header.Add("Authorization", "Bearer " + fs.Config.FACEIT_API_KEY)
	response, err := client.Fire(req)
	body, err := ioutil.ReadAll(response.Body)

	log.Println("Match End")
	log.Println(string(body))

	var stats faceit.Match
	_ = json.Unmarshal(body, &stats)

	for _, s := range stats.Rounds {
		for _, a := range s.Teams {
			for _, d := range a.Players {
				if d.ID == "0d94613d-b736-46ba-b8cd-d2159ddad705" || d.ID == "b26df7d4-8517-4ec6-ab58-708487e5fe60" || d.ID == "b0a57a5a-2f7a-481c-aaa8-8013a83378e3" {

					var outcome = "Victory"

					if d.Stats.Result == "0" {
						outcome = "Defeat"
					}

					*messages = append(*messages, helpers.NewEmbed().
						SetTitle("Match ended for "+d.Nickname).
						SetDescription(outcome + " on " + stats.Rounds[0].MatchStats.Map + " [" +
							stats.Rounds[0].MatchStats.Score + "]").
						AddField("Kills", d.Stats.Kills, true).
						AddField("Assists", d.Stats.Assists, true).
						AddField("Deaths", d.Stats.Deaths, true).
						AddField("K/D Ratio", d.Stats.KD, true).
						AddField("K/R Ratio", d.Stats.KR, true))
				}
			}
		}
	}

	if err != nil {
		log.Println(err)
	}
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
				if d.ID == "0d94613d-b736-46ba-b8cd-d2159ddad705" || d.ID == "b26df7d4-8517-4ec6-ab58-708487e5fe60" || d.ID == "b0a57a5a-2f7a-481c-aaa8-8013a83378e3" {
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
	req, err := http.NewRequest("GET", api.GetFaceitMatch(webhook.Payload.MatchID), nil)
	req.Header.Add("Authorization", "Bearer " + fs.Config.FACEIT_API_KEY)
	response, err := client.Fire(req)
	body, err := ioutil.ReadAll(response.Body)

	log.Println("Match Ready")
	log.Println(string(body))

	var stats faceit.Match
	_ = json.Unmarshal(body, &stats)

	for _, s := range stats.Rounds {
		for _, a := range s.Teams {
			for _, d := range a.Players {
				if d.ID == "0d94613d-b736-46ba-b8cd-d2159ddad705" || d.ID == "b26df7d4-8517-4ec6-ab58-708487e5fe60" || d.ID == "b0a57a5a-2f7a-481c-aaa8-8013a83378e3" {
					*messages = append(*messages, helpers.NewEmbed().
						SetTitle("Match ready for "+d.Nickname))
				}
			}
		}
	}

	if err != nil {
		log.Println(err)
	}
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
				if d.ID == "0d94613d-b736-46ba-b8cd-d2159ddad705" || d.ID == "b26df7d4-8517-4ec6-ab58-708487e5fe60" || d.ID == "b0a57a5a-2f7a-481c-aaa8-8013a83378e3" {
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