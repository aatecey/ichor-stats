package api

import (
	"bytes"
	client "ichor-stats/src/package/http"
	"io/ioutil"
	"log"
	"net/http"
)

var ENDPOINT = "http://40.121.59.132:5000"

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
	return ENDPOINT + "/match/stats"
}

func GetAllSinglePlayerStatsEndpoint() string {
	return ENDPOINT + "/player/stats"
}

func GetLifetimePlayerStatsEndpoint() string {
	return ENDPOINT + "/player/lifetime"
}