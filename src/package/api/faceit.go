package api

import (
	"fmt"
	"ichor-stats/src/app/services/config"
	client "ichor-stats/src/package/http"
	"io/ioutil"
	"log"
	"net/http"
)

func FaceitRequest(apiUrl string) []byte  {
	var bearer = "Bearer " + config.GetConfig().FACEIT_API_KEY
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Println("hereeeee")
		log.Println(err)
		return nil
	}

	req.Header.Add("Authorization", bearer)
	response, err := client.Fire(req)

	log.Println("Client Fired Completed + " + apiUrl)

	if err != nil {
		log.Println("In here", err)
		//return nil
	}

	body, err := ioutil.ReadAll(response.Body)
	return body
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
