package discord

import (
	"github.com/bwmarrin/discordgo"
	"ichor-stats/src/app/models/faceit"
	"ichor-stats/src/app/services/config"
	"ichor-stats/src/app/services/discord/helpers"
	"ichor-stats/src/package/api"
	"log"
	"strconv"
	"strings"
)

func EmbeddedStats(stats faceit.Stats, user faceit.User) *helpers.Embed {
	kills, assists, deaths := DetermineTotalStats(stats, user)

	return helpers.NewEmbed().
		SetTitle(user.Games.CSGO.Name).
		AddField("ELO", strconv.Itoa(user.Games.CSGO.ELO), true).
		AddField("Skill Level", strconv.Itoa(user.Games.CSGO.SkillLevel), true).
		AddField("Avg. K/D Ratio", stats.Lifetime.AverageKD, false).
		AddField("Avg. Headshots %", stats.Lifetime.AverageHeadshots, false).
		AddField("Total Kills", kills, true).
		AddField("Total Assists", assists, true).
		AddField("Total Deaths", deaths, true)
}


func EmbeddedStreak(stats faceit.Stats, user faceit.User) *helpers.Embed {
	var resultsArray []string
	for _, result := range stats.Lifetime.RecentResults {
		if result == "0" {
			resultsArray =append(resultsArray, "L")
		} else {
			resultsArray =append(resultsArray, "W")
		}
	}

	return helpers.NewEmbed().
		SetTitle(user.Games.CSGO.Name).
		AddField("Recent Results (Most recent on right)", strings.Join(resultsArray, ", "), false).
		AddField("Current Win Streak", stats.Lifetime.CurrentWinStreak, false)
}

func EmbeddedGreen() *helpers.Embed {
	return helpers.NewEmbed().
		SetTitle("Lmao green fkn noob.")
}

func EmbeddedMapStats(stats faceit.Stats, user faceit.User, gameMap string) *helpers.Embed {
	for _, result := range stats.Segment {
		if strings.HasSuffix(result.CsMap, gameMap) {
			return helpers.NewEmbed().
				SetTitle("Map statistics for " + user.Games.CSGO.Name + " on " + result.CsMap).
				AddField("Kills", result.LifetimeMapStats.Kills, true).
				AddField("Assists", result.LifetimeMapStats.Assists, true).
				AddField("Deaths", result.LifetimeMapStats.Deaths, true).
				AddField("Triple Kills", result.LifetimeMapStats.TripleKills, true).
				AddField("Quadro Kills", result.LifetimeMapStats.QuadroKills, true).
				AddField("Penta Kills", result.LifetimeMapStats.PentaKills, true).
				AddField("Avg. K/D Ratio", result.LifetimeMapStats.AverageKD, true).
				AddField("Win Rate (%)", result.LifetimeMapStats.WinRate, true)
		}
	}

	return nil
}

func EmbeddedLastFiveStats(stats faceit.Stats, user faceit.User, numberOfMatches string, requesterID string, session *discordgo.Session) *helpers.Embed {
	var matchHistory faceit.Matches
	matchesErr := GetFaceitStats(api.GetFaceitPlayerMatchHistory(requesterID)).Decode(&matchHistory)
	if matchesErr != nil {
		log.Println(matchesErr)
	}

	matchNumber, err := strconv.Atoi(numberOfMatches)
	if err != nil {
		log.Println(err)
	}

	if matchNumber > 20 {
		matchNumber = 20
	}

	var listOfMatches = make([]MatchDetails, matchNumber)

	for i := 0; i < matchNumber; i++ {
		var match faceit.Match
		matchErr := GetFaceitStats(api.GetFaceitMatchDetails(matchHistory.MatchItem[i].MatchId)).Decode(&match)
		if matchErr != nil {
			log.Println(matchesErr)
		}

		log.Println(match.Rounds[0].MatchStats.Score)

		listOfMatches[i].MapName = match.Rounds[0].MatchStats.Map
		listOfMatches[i].Score = match.Rounds[0].MatchStats.Score

		for _, team := range match.Rounds[0].Teams {
			for _, player := range team.Players {
				if player.ID == requesterID {
					listOfMatches[i].Kills = player.Stats.Kills
					listOfMatches[i].Assists = player.Stats.Assists
					listOfMatches[i].Deaths = player.Stats.Deaths
					listOfMatches[i].Result = player.Stats.Result
					break
				}
			}
		}
	}

	timeSinceGame := 1

	for _, game := range listOfMatches {
		outcome := "Win"

		if game.Result == "0" {
			outcome = "Loss"
		}

		embeddedMsg := helpers.NewEmbed().
			SetTitle(user.Games.CSGO.Name + " played " + game.MapName + " " + strconv.Itoa(timeSinceGame) + " game(s) ago:").
			AddField("Kills", game.Kills, true).
			AddField("Assists", game.Assists, true).
			AddField("Deaths", game.Deaths, true).
			AddField("Result", outcome + " [" + game.Score + "]", false)

		timeSinceGame++

		if embeddedMsg != nil {
			_, err := session.ChannelMessageSendEmbed(config.GetConfig().CHANNEL_ID, embeddedMsg.MessageEmbed)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return nil
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


type MatchDetails struct {
	MapName	string
	Kills string
	Deaths string
	Assists string
	Result string
	Score string
}