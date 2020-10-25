package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"ichor-stats/src/app/models/config"
	"ichor-stats/src/app/models/faceit"
	"ichor-stats/src/app/services/discord/helpers"
	"ichor-stats/src/app/services/firebase"
	"ichor-stats/src/package/api"
	"log"
	"strconv"
	"strings"
)

type ServiceDiscord struct {
	Config config.Configuration
	Discord discordgo.Session
}

func NewDiscordService(config *config.Configuration) ServiceDiscord {
	dc, err := discordgo.New("Bot " + config.DISCORD_BOT_ID)
	if err != nil {
		fmt.Println(err)
	}

	return ServiceDiscord{
		Config: *config,
		Discord: *dc,
	}
}

func (ds *ServiceDiscord) SendMessage(message string) {
	_, err := ds.Discord.ChannelMessageSend(ds.Config.CHANNEL_ID, message)
	if err != nil {
		fmt.Println(err)
	}
}

func HandleParameterisedCommand(requesterId string, command []string, user faceit.User, messages *[]*helpers.Embed) {
	switch trimLeftChar(command[0]) {
	case "map":
		helpers.EmbeddedMapStats(requesterId, user, command[1], messages)
	case "last":
		firebase.EmbeddedLastMatchStats(requesterId, user, command[1], messages)
	case "totals":
		firebase.EmbeddedLastMatchTotals(requesterId, user, command[1], messages)
	}
}

func HandleCommand(requesterId string, command string, user faceit.User, messages *[]*helpers.Embed) {
	switch trimLeftChar(command) {
	case "stats":
		EmbeddedStats(requesterId, user, messages)
	case "streak":
		EmbeddedStreak(requesterId, user, messages)
	case "recent":
		EmbeddedStreak(requesterId, user, messages)
	case "green":
		EmbeddedGreen(messages)
	}
}

func EmbeddedStats(requesterId string, user faceit.User, messages *[]*helpers.Embed) {
	var stats faceit.Stats
	_ = api.FaceitRequest(api.GetFaceitPlayerCsgoStats(requesterId)).Decode(&stats)
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
	_ = api.FaceitRequest(api.GetFaceitPlayerCsgoStats(requesterId)).Decode(&stats)

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

func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}