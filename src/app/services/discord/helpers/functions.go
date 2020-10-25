package helpers

import (
	"ichor-stats/src/app/models/faceit"
	"ichor-stats/src/package/api"
	"log"
	"strconv"
	"strings"
)

func EmbeddedMapStats(requesterId string, user faceit.User, gameMap string, messages *[]*Embed) {
	var stats faceit.Stats
	_ = api.FaceitRequest(api.GetFaceitPlayerCsgoStats(requesterId)).Decode(&stats)

	for _, result := range stats.Segment {
		if strings.HasSuffix(result.CsMap, gameMap) {
			*messages = append(*messages, NewEmbed().
				SetTitle("Map statistics for "+user.Games.CSGO.Name+" on "+result.CsMap).
				AddField("Kills", result.LifetimeMapStats.Kills, true).
				AddField("Assists", result.LifetimeMapStats.Assists, true).
				AddField("Deaths", result.LifetimeMapStats.Deaths, true).
				AddField("Triple Kills", result.LifetimeMapStats.TripleKills, true).
				AddField("Quadro Kills", result.LifetimeMapStats.QuadroKills, true).
				AddField("Penta Kills", result.LifetimeMapStats.PentaKills, true).
				AddField("Avg. K/D Ratio", result.LifetimeMapStats.AverageKD, true).
				AddField("Win Rate (%)", result.LifetimeMapStats.WinRate, true))
			return
		}
	}
}

func EmbeddedLastMatchStats(requesterID string, user faceit.User, numberOfMatches string, messages *[]*Embed) {
	matchHistory, totalMatches := GetMatchHistory(numberOfMatches, requesterID)

	for i := 0; i < totalMatches; i++ {
		var match faceit.Match
		matchErr := api.FaceitRequest(api.GetFaceitMatchDetails(matchHistory.MatchItem[i].MatchId)).Decode(&match)
		if matchErr != nil {
			log.Println(matchErr)
		}

		var stats = GetPlayerDetailsFromMatch(match, requesterID)

		outcome := "Win"
		if stats.Result == "0" {
			outcome = "Loss"
		}

		*messages = append(*messages, NewEmbed().
			SetTitle(user.Games.CSGO.Name+" played "+match.Rounds[0].MatchStats.Map+" "+strconv.Itoa(i+1)+" game(s) ago:").
			AddField("Kills", stats.Kills, true).
			AddField("Assists", stats.Assists, true).
			AddField("Deaths", stats.Deaths, true).
			AddField("Result", outcome+" ["+match.Rounds[0].MatchStats.Score+"]", false))
	}
}

func EmbeddedLastMatchTotals(requesterID string, user faceit.User, numberOfMatches string, messages *[]*Embed) {
	matchHistory, totalMatches := GetMatchHistory(numberOfMatches, requesterID)

	totalKills := 0
	totalAssists := 0
	totalDeaths := 0
	totalWins := 0

	for i := 0; i < totalMatches; i++ {
		var match faceit.Match
		err := api.FaceitRequest(api.GetFaceitMatchDetails(matchHistory.MatchItem[i].MatchId)).Decode(&match)
		var stats = GetPlayerDetailsFromMatch(match, requesterID)

		gameKills, err := strconv.Atoi(stats.Kills)
		gameAssists, err := strconv.Atoi(stats.Assists)
		gameDeaths, err := strconv.Atoi(stats.Deaths)

		totalKills = totalKills + gameKills
		totalAssists = totalAssists + gameAssists
		totalDeaths = totalDeaths + gameDeaths

		if stats.Result == "1" {
			totalWins++
		}

		if err != nil {
			log.Println(err)
		}
	}

	*messages = append(*messages, NewEmbed().
		SetTitle(user.Games.CSGO.Name+" stats for the last "+strconv.Itoa(totalMatches)+" games:").
		AddField("Kills", strconv.Itoa(totalKills), true).
		AddField("Assists", strconv.Itoa(totalAssists), true).
		AddField("Deaths", strconv.Itoa(totalDeaths), true).
		AddField("Wins", strconv.Itoa(totalWins), false))
}

func GetMatchHistory(numberOfMatches string, requesterID string) (history faceit.Matches, totalMatches int) {
	var matchHistory faceit.Matches
	matchesErr := api.FaceitRequest(api.GetFaceitPlayerMatchHistory(requesterID)).Decode(&matchHistory)
	userRequestedMatches, err := strconv.Atoi(numberOfMatches)
	var matchLimit, err2 = strconv.Atoi(numberOfMatches)

	if err != nil || err2 != nil || matchesErr != nil {
		log.Println(err)
	}

	if userRequestedMatches > len(matchHistory.MatchItem) {
		matchLimit = len(matchHistory.MatchItem)
	}

	return matchHistory, matchLimit
}

func GetPlayerDetailsFromMatch(match faceit.Match, requesterID string) (stats faceit.PlayerStats) {
	for _, team := range match.Rounds[0].Teams {
		for _, player := range team.Players {
			if player.ID == requesterID {
				return faceit.PlayerStats {
					Kills: player.Stats.Kills,
					Assists: player.Stats.Assists,
					Deaths: player.Stats.Deaths,
					Result: player.Stats.Result,
					KD: player.Stats.KD,
					KR: player.Stats.KR,
					Headshots: player.Stats.Headshots,
					HeadshotPercentage: player.Stats.HeadshotPercentage,
					Pentas: player.Stats.Pentas,
					Quads: player.Stats.Quads,
					Triples: player.Stats.Triples,
					MVPs: player.Stats.MVPs,
				}
			}
		}
	}

	return faceit.PlayerStats{}
}