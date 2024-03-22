package calls

import (
	"encoding/json"
	"ichor-stats/src/package/api"
	"ichor-stats/src/package/discord"
	"log"
	"strconv"
	"strings"
)

type ResponseJson struct {
	TotalMatches int            `json:"Total Matches"`
	TotalKills	 int			`json:"Total Kills"`
	TotalDeaths	 int			`json:"Total Deaths"`
	KillsArray   []int          `json:"KillsArray"`
	DeathsArray  []int          `json:"DeathsArray"`
	Assists      int            `json:"Assists"`
	Triples      int            `json:"Triples"`
	Quads        int            `json:"Quads"`
	Pentas       int            `json:"Pentas"`
	MVPS         int            `json:"MVPS"`
	Wins         int            `json:"Wins"`
	Losses       int            `json:"Losses"`
	MapStats     map[string]int `json:"MapStats"`
}

type MatchResponseJson struct {
	Matches []Match `json:"Matches"`
}

type Match struct {
	Map     string `json:"Map"`
	Kills   string `json:"Kills"`
	Assists string `json:"Assists"`
	Deaths  string `json:"Deaths"`
	Result  string `json:"Result"`
	Score   string `json:"Score"`
}

type LifetimeResponseJson struct {
	SkillLevel       string `json:"skill_level"`
	ELO              string `json:"elo"`
	AverageHeadshots string `json:"Average Headshots %"`
	AverageKD        string `json:"Average K/D Ratio"`
	CurrentWinStreak string `json:"Current Win Streak"`
	TotalKills       string `json:"Total Kills"`
	TotalDeaths      string `json:"Total Deaths"`
	TotalAssists     string `json:"Total Assists"`
	LifetimeMapStats map[string]LifetimeMapStats
}

type LifetimeMapStats struct {
	Assists     string `json:"Assists"`
	Kills       string `json:"Kills"`
	Deaths      string `json:"Deaths"`
	WinRate     string `json:"Win Rate %"`
	AverageKD   string `json:"Average K/D Ratio"`
	TripleKills string `json:"Triple Kills"`
	QuadroKills string `json:"Quadro Kills"`
	PentaKills  string `json:"Penta Kills"`
}

func MapStats(playerName string, gameMap string, messages *[]*discord.Embed) {
	var stats LifetimeResponseJson
	_ = json.Unmarshal(api.ApiRequest(api.GetLifetimePlayerStatsEndpoint(), "", playerName, "true"), &stats)

	for mapName, mapStats := range stats.LifetimeMapStats {
		if strings.HasSuffix(mapName, gameMap) {
			*messages = append(*messages, discord.NewEmbed().
				SetTitle("Map statistics for "+playerName+" on "+mapName).
				AddField("Kills", mapStats.Kills, true).
				AddField("Assists", mapStats.Assists, true).
				AddField("Deaths", mapStats.Deaths, true).
				AddField("Triple Kills", mapStats.TripleKills, true).
				AddField("Quadro Kills", mapStats.QuadroKills, true).
				AddField("Penta Kills", mapStats.PentaKills, true).
				AddField("Avg. K/D Ratio", mapStats.AverageKD, true).
				AddField("Win Rate (%)", mapStats.WinRate, true))
			return
		}
	}
}

func LastMatchStats(playerName string, numberOfMatches string, oldestMatchFirst string, messages *[]*discord.Embed) {
	var stats MatchResponseJson
	_ = json.Unmarshal(api.ApiRequest(api.GetMatchStatsForPlayerEndpoint(), numberOfMatches, playerName, oldestMatchFirst), &stats)

	for i := 0; i < len(stats.Matches); i++ {
		*messages = append(*messages, discord.NewEmbed().
			SetTitle(playerName+" played "+stats.Matches[i].Map+" "+strconv.Itoa(i+1)+" game(s) ago:").
			AddField("Kills", stats.Matches[i].Kills, true).
			AddField("Assists", stats.Matches[i].Assists, true).
			AddField("Deaths", stats.Matches[i].Deaths, true).
			AddField("Result", stats.Matches[i].Result+" ["+stats.Matches[i].Score+"]", false))
	}
}

func LastMatchTotals(playerName string, numberOfMatches string, oldestMatchFirst string, messages *[]*discord.Embed) {
	var stats map[string]ResponseJson
	_ = json.Unmarshal(api.ApiRequest(api.GetAllSinglePlayerStatsEndpoint(), numberOfMatches, playerName, oldestMatchFirst), &stats)

	*messages = append(*messages, discord.NewEmbed().
		SetTitle(playerName + " stats for the last "+strconv.Itoa(stats[playerName].TotalMatches) + " games:").
		AddField("Kills", strconv.Itoa(stats[playerName].TotalKills), true).
		AddField("Assists", strconv.Itoa(stats[playerName].Assists), true).
		AddField("Deaths", strconv.Itoa(stats[playerName].TotalDeaths), true).
		AddField("Wins", strconv.Itoa(stats[playerName].Wins), false))
}

func Stats(playerName string, messages *[]*discord.Embed) {
	var stats LifetimeResponseJson
	_ = json.Unmarshal(api.ApiRequest(api.GetLifetimePlayerStatsEndpoint(), "", playerName, "true"), &stats)

	*messages = append(*messages, discord.NewEmbed().
		SetTitle(playerName).
		AddField("ELO", stats.ELO, true).
		AddField("Skill Level", stats.SkillLevel, true).
		AddField("Avg. K/D Ratio", stats.AverageKD, false).
		AddField("Avg. Headshots %", stats.AverageHeadshots, false).
		AddField("Total Kills", stats.TotalKills, true).
		AddField("Total Assists", stats.TotalAssists, true).
		AddField("Total Deaths", stats.TotalDeaths, true))

	log.Println("Adding message: " + strconv.Itoa(len(*messages)))
}

func Streak(playerName string, oldestMatchFirst string, messages *[]*discord.Embed) {
	var stats MatchResponseJson
	_ = json.Unmarshal(api.ApiRequest(api.GetMatchStatsForPlayerEndpoint(), "20", playerName, oldestMatchFirst), &stats)

	var resultsArray []string
	var winStreak = 0
	winStreakEnd := false

	for i := 0; i < len(stats.Matches); i++ {
		if stats.Matches[i].Result == "Win" {
			if len(resultsArray) < 5 {
				resultsArray = append(resultsArray, "W")
			}
			if !winStreakEnd {
				winStreak++
			}
		} else {
			winStreakEnd = true
			if len(resultsArray) < 5 {
				resultsArray = append(resultsArray, "L")
			}
		}
	}

	*messages = append(*messages, discord.NewEmbed().
		SetTitle(playerName).
		AddField("Recent Results (Most recent on left)", strings.Join(resultsArray, ", "), false).
		AddField("Current Win Streak", strconv.Itoa(winStreak), false))
}

func Green(messages *[]*discord.Embed) {
	*messages = append(*messages, discord.NewEmbed().
		SetTitle("Green sucks at Smash Bros too."))
}
