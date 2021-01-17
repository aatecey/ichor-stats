package api

import (
	"bytes"
	client "ichor-stats/src/package/http"
	"io/ioutil"
	"log"
	"net/http"
)

func ApiRequest(apiUrl string, numberOfMatches string, playerName string, oldestMatchFirst string) []byte {
	var jsonStr = []byte(`{"matchCount":"` + numberOfMatches + `","oldestMatchFirst":` + oldestMatchFirst + `,"name":"` + playerName + `"}`)

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println("Error when forming HTTP request to fire against FaceIt - ", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	response, err := client.Fire(req)
	defer response.Body.Close()
	if err != nil {
		log.Println("Error when firing request against FaceIt - ", err)
	}

	body, err := ioutil.ReadAll(response.Body)
	return body
}

func GetMatchStatsForPlayerEndpoint() string {
	return "http://127.0.0.1:5000/match/stats"
}

func GetAllSinglePlayerStatsEndpoint() string {
	return "http://127.0.0.1:5000/player/stats"
}

func GetLifetimePlayerStatsEndpoint() string {
	return "http://127.0.0.1:5000/player/lifetime"
}