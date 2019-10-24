package faceit

import (
	"encoding/json"
	"ichor-stats/src/app/models/config"
	"ichor-stats/src/app/models/faceit"
	"ichor-stats/src/app/services/discord"
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

func (fs *ServiceFaceit) MatchEnd(webhook faceit.Webhook) error {

	log.Println(webhook.Payload.MatchID)
	req, err := http.NewRequest("GET", api.GetFaceitMatch(webhook.Payload.MatchID), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	var bearer = "Bearer " + fs.Config.FACEIT_API_KEY
	req.Header.Add("Authorization", bearer)

	response, err := client.Fire(req)
	if err != nil {
		log.Println(err)
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	bodyString := string(body)
	log.Println(bodyString)

	var stats faceit.Match
	err = json.NewDecoder(response.Body).Decode(&stats)
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