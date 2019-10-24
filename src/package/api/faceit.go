package api

import "fmt"

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
