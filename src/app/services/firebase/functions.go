package firebase

import (
	"ichor-stats/src/app/models/faceit"
	"ichor-stats/src/app/services/discord/helpers"
	"log"
	"strconv"
)

func EmbeddedLastMatchStats(requesterID string, user faceit.User, numberOfMatches string, messages *[]*helpers.Embed) {
	var matchesFromDb = GetMatchStats(numberOfMatches, requesterID)

	for i := 0; i < len(matchesFromDb); i++ {
		*messages = append(*messages, helpers.NewEmbed().
			SetTitle(user.Games.CSGO.Name+" played "+matchesFromDb[i].Map+" "+strconv.Itoa(i+1)+" game(s) ago:").
			AddField("Kills", matchesFromDb[i].Kills, true).
			AddField("Assists", matchesFromDb[i].Assists, true).
			AddField("Deaths", matchesFromDb[i].Deaths, true).
			AddField("Result", matchesFromDb[i].Result+" ["+matchesFromDb[i].Score+"]", false))
	}
}

func EmbeddedLastMatchTotals(requesterID string, user faceit.User, numberOfMatches string, messages *[]*helpers.Embed) {
	var matchesFromDb = GetMatchStats(numberOfMatches, requesterID)

	totalKills := 0
	totalAssists := 0
	totalDeaths := 0
	totalWins := 0

	for i := 0; i < len(matchesFromDb); i++ {
		gameKills, err := strconv.Atoi(matchesFromDb[i].Kills)
		gameAssists, err := strconv.Atoi(matchesFromDb[i].Assists)
		gameDeaths, err := strconv.Atoi(matchesFromDb[i].Deaths)

		totalKills = totalKills + gameKills
		totalAssists = totalAssists + gameAssists
		totalDeaths = totalDeaths + gameDeaths

		if matchesFromDb[i].Result == "Win" {
			totalWins++
		}

		if err != nil {
			log.Println(err)
		}
	}

	*messages = append(*messages, helpers.NewEmbed().
		SetTitle(user.Games.CSGO.Name+" stats for the last "+strconv.Itoa(len(matchesFromDb))+" games:").
		AddField("Kills", strconv.Itoa(totalKills), true).
		AddField("Assists", strconv.Itoa(totalAssists), true).
		AddField("Deaths", strconv.Itoa(totalDeaths), true).
		AddField("Wins", strconv.Itoa(totalWins), false))
}