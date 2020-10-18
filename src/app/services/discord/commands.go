package discord

import (
	"ichor-stats/src/app/models/faceit"
	"ichor-stats/src/app/services/discord/helpers"
	"ichor-stats/src/package/api"
	"log"
	"strconv"
	"strings"
)

func EmbeddedStats(requesterId string, user faceit.User, messages *[]*helpers.Embed) {
	var stats faceit.Stats
	_ = GetFaceitStats(api.GetFaceitPlayerCsgoStats(requesterId)).Decode(&stats)
	kills, assists, deaths := DetermineTotalStats(stats, user)

	*messages = append(*messages, helpers.NewEmbed().
		SetTitle(user.Games.CSGO.Name).
		AddField("ELO", strconv.Itoa(user.Games.CSGO.ELO), true).
		AddField("Skill Level", strconv.Itoa(user.Games.CSGO.SkillLevel), true).
		AddField("Avg. K/D Ratio", stats.Lifetime.AverageKD, false).
		AddField("Avg. Headshots %", stats.Lifetime.AverageHeadshots, false).
		AddField("Total Kills", kills, true).
		AddField("Total Assists", assists, true).
		AddField("Total Deaths", deaths, true))

	log.Println("Adding message: " + strconv.Itoa(len(*messages)))
}


func EmbeddedStreak(requesterId string, user faceit.User, messages *[]*helpers.Embed) {
	var stats faceit.Stats
	_ = GetFaceitStats(api.GetFaceitPlayerCsgoStats(requesterId)).Decode(&stats)

	var resultsArray []string
	for _, result := range stats.Lifetime.RecentResults {
		if result == "0" {
			resultsArray = append(resultsArray, "L")
		} else {
			resultsArray = append(resultsArray, "W")
		}
	}

	*messages = append(*messages, helpers.NewEmbed().
		SetTitle(user.Games.CSGO.Name).
		AddField("Recent Results (Most recent on right)", strings.Join(resultsArray, ", "), false).
		AddField("Current Win Streak", stats.Lifetime.CurrentWinStreak, false))
}

func EmbeddedGreen(messages *[]*helpers.Embed) {
	*messages = append(*messages, helpers.NewEmbed().
		SetTitle("Green stop holding connector"))
}

func EmbeddedMapStats(requesterId string, user faceit.User, gameMap string, messages *[]*helpers.Embed) {
	var stats faceit.Stats
	_ = GetFaceitStats(api.GetFaceitPlayerCsgoStats(requesterId)).Decode(&stats)

	for _, result := range stats.Segment {
		if strings.HasSuffix(result.CsMap, gameMap) {
			*messages = append(*messages, helpers.NewEmbed().
				SetTitle("Map statistics for " + user.Games.CSGO.Name + " on " + result.CsMap).
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

func EmbeddedLastMatchStats(requesterID string, user faceit.User, numberOfMatches string, messages *[]*helpers.Embed) {
	matchHistory, totalMatches := GetMatchHistory(numberOfMatches, requesterID)

	for i := 0; i < totalMatches; i++ {
		var match faceit.Match
		matchErr := GetFaceitStats(api.GetFaceitMatchDetails(matchHistory.MatchItem[i].MatchId)).Decode(&match)
		if matchErr != nil {
			log.Println(matchErr)
		}

		var stats = GetPlayerDetailsFromMatch(match, requesterID)

		outcome := "Win"
		if stats.Result == "0" {
			outcome = "Loss"
		}

		*messages = append(*messages, helpers.NewEmbed().
			SetTitle(user.Games.CSGO.Name + " played " + match.Rounds[0].MatchStats.Map + " " + strconv.Itoa(i + 1) + " game(s) ago:").
			AddField("Kills", stats.Kills, true).
			AddField("Assists", stats.Assists, true).
			AddField("Deaths", stats.Deaths, true).
			AddField("Result", outcome + " [" + match.Rounds[0].MatchStats.Score + "]", false))
	}
}

func EmbeddedLastMatchTotals(requesterID string, user faceit.User, numberOfMatches string, messages *[]*helpers.Embed) {
	matchHistory, totalMatches := GetMatchHistory(numberOfMatches, requesterID)

	totalKills := 0
	totalAssists := 0
	totalDeaths := 0
	totalWins := 0

	for i := 0; i < totalMatches; i++ {
		var match faceit.Match
		err := GetFaceitStats(api.GetFaceitMatchDetails(matchHistory.MatchItem[i].MatchId)).Decode(&match)
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

	*messages = append(*messages, helpers.NewEmbed().
		SetTitle(user.Games.CSGO.Name+" stats for the last " + strconv.Itoa(totalMatches) + " games:").
		AddField("Kills", strconv.Itoa(totalKills), true).
		AddField("Assists", strconv.Itoa(totalAssists), true).
		AddField("Deaths", strconv.Itoa(totalDeaths), true).
		AddField("Wins", strconv.Itoa(totalWins), false))
}

func GetMatchHistory(numberOfMatches string, requesterID string) (history faceit.Matches, totalMatches int) {
	var matchHistory faceit.Matches
	matchesErr := GetFaceitStats(api.GetFaceitPlayerMatchHistory(requesterID)).Decode(&matchHistory)
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
				}
			}
		}
	}

	return faceit.PlayerStats{}
}

func DetermineTotalStats(stats faceit.Stats, user faceit.User) (string, string, string) {
	totalKills := 0
	totalDeaths := 0
	totalAssists := 0

	for _, result := range stats.Segment {
		mapKills, _ := strconv.Atoi(result.LifetimeMapStats.Kills)
		mapDeaths, _ := strconv.Atoi(result.LifetimeMapStats.Deaths)
		mapAssists, _ := strconv.Atoi(result.LifetimeMapStats.Assists)

		totalKills = totalKills + mapKills
		totalDeaths = totalDeaths + mapDeaths
		totalAssists = totalAssists + mapAssists
	}

	return strconv.Itoa(totalKills), strconv.Itoa(totalAssists), strconv.Itoa(totalDeaths)
}