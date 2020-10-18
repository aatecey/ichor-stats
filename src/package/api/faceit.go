package api

import (
	"encoding/json"
	"fmt"
	"ichor-stats/src/app/services/config"
	client "ichor-stats/src/package/http"
	"log"
	"net/http"
)

func FaceitRequest(apiUrl string) *json.Decoder  {
	var bearer = "Bearer " + config.GetConfig().FACEIT_API_KEY
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Println(err)
		return nil
	}

	req.Header.Add("Authorization", bearer)
	response, err := client.Fire(req)
	if err != nil {
		log.Println(err)
		return nil
	}

	return json.NewDecoder(response.Body)
}

func GetFaceitPlayerCsgoStats(playerId string) string {
	return fmt.Sprintf("https://open.faceit.com/data/v4/players/%s/stats/csgo", playerId)
}

func GetFaceitPlayerStats(playerId string) string {
	return fmt.Sprintf("https://open.faceit.com/data/v4/players/%s", playerId)
}

func GetFaceitMatch(playerId string) string {
	return fmt.Sprintf("https://open.faceit.com/data/v4/matches/%s/stats", playerId)
}

func GetFaceitPlayerMatchHistory(playerId string) string {
	return fmt.Sprintf("https://open.faceit.com/data/v4/players/%s/history", playerId)
}

func GetFaceitMatchDetails(matchId string) string {
	return fmt.Sprintf("https://open.faceit.com/data/v4/matches/%s/stats", matchId)
}
