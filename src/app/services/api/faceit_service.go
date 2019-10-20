package api

import (
	"encoding/json"
	"ichor-stats/src/app/models/config"
	"ichor-stats/src/app/models/faceit"
	"ichor-stats/src/app/services/discord"
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

func (fs *ServiceFaceit) MatchEnd(webhook faceit.Webhook) error {
	url := "https://open.faceit.com/data/v4/matches/" + webhook.Payload.MatchID + "/stats"

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + fs.Config.FACEIT_API_KEY

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	var stats faceit.Match
	err = json.NewDecoder(resp.Body).Decode(&stats)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(&stats)
	log.Println(stats)
	for _, s := range stats.Rounds {
		for _, a := range s.Teams {
			log.Println(a.MatchStats.Map)
			log.Println(a.MatchStats.Score)
			for _, d := range a.Players {
				if d.ID == "0d94613d-b736-46ba-b8cd-d2159ddad705" || d.ID == "b26df7d4-8517-4ec6-ab58-708487e5fe60" || d.ID == "b0a57a5a-2f7a-481c-aaa8-8013a83378e3" {
					log.Println(d.Stats.Kills)
					log.Println(d.Stats.Assists)
					log.Println(d.Stats.Deaths)
					log.Println(d.Stats.KD)
					log.Println(d.Stats.KR)
				}
			}
		}
	}

	fs.DiscordService.SendMessage("Match ended")
	return nil
}
