package firebase

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/option"
	"ichor-stats/src/app/models/database"
	"ichor-stats/src/app/models/faceit"
	"ichor-stats/src/app/models/players"
	"ichor-stats/src/app/services/discord/helpers"
	"ichor-stats/src/package/api"
	"log"
	"strconv"
)

var client *db.Client

func Init() {
	opt := option.WithCredentialsFile("./src/build/firebase-creds.json")
	config := &firebase.Config{
		ProjectID: "ichor-stats-db",
		DatabaseURL: "https://ichor-stats-db.firebaseio.com",
	}

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, _ = app.Database(context.Background())
}

func Setup() {
	for player := range players.Players {
		var matches = make([]database.Match, 0)

		var playerDetails players.PlayerDetails
		_ = mapstructure.Decode(players.Players[player], &playerDetails)

		matchHistory, _ := helpers.GetMatchHistory("100", player)

		for _, match := range matchHistory.MatchItem {
			var matchDetails faceit.Match
			matchErr := api.FaceitRequest(api.GetFaceitMatchDetails(match.MatchId)).Decode(&matchDetails)
			if matchErr != nil {
				log.Println(matchErr)
			}

			var stats = helpers.GetPlayerDetailsFromMatch(matchDetails, player)
			var mapName = matchDetails.Rounds[0].MatchStats.Map

			outcome := "Win"
			if stats.Result == "0" {
				outcome = "Loss"
			}

			match := database.Match{
				ID:                 match.MatchId,
				Map:                mapName,
				Result:             outcome,
				Score:              matchDetails.Rounds[0].MatchStats.Score,
				Kills:              stats.Kills,
				Assists:            stats.Assists,
				Deaths:             stats.Deaths,
				Headshots:          stats.Headshots,
				HeadshotPercentage: stats.HeadshotPercentage,
				Pentas:             stats.Pentas,
				Quads:              stats.Quads,
				Triples:            stats.Triples,
				KillDeathRatio:     stats.KD,
				KillRoundRatio:     stats.KR,
				MVPs:               stats.MVPs,
			}

			matches = append(matches, match)
		}

		player := database.Player{
			Matches: matches,
		}

		if err := client.NewRef(playerDetails.Name).Set(context.Background(), player); err != nil {
			log.Fatal(err)
		}
	}
}

func SaveMatch(playerFromMatch faceit.Players, teams faceit.Teams, matchStats faceit.Match, matchId string) {
	matchesFromDb := GetMatchStats("10000", playerFromMatch.ID)

	var outcome = "Win"

	if playerFromMatch.Stats.Result == "0" {
		outcome = "Loss"
	}

	match := database.Match{
		ID:                 matchId,
		Map:                teams.MatchStats.Map,
		Result:             outcome,
		Score:              matchStats.Rounds[0].MatchStats.Score,
		Kills:              playerFromMatch.Stats.Kills,
		Assists:            playerFromMatch.Stats.Assists,
		Deaths:             playerFromMatch.Stats.Deaths,
		Headshots:          playerFromMatch.Stats.Headshots,
		HeadshotPercentage: playerFromMatch.Stats.HeadshotPercentage,
		Pentas:             playerFromMatch.Stats.Pentas,
		Quads:              playerFromMatch.Stats.Quads,
		Triples:            playerFromMatch.Stats.Triples,
		KillDeathRatio:     playerFromMatch.Stats.KD,
		KillRoundRatio:     playerFromMatch.Stats.KR,
		MVPs:               playerFromMatch.Stats.MVPs,
	}

	matchesFromDb = append([]database.Match{match}, matchesFromDb...)

	player := database.Player{
		Matches: matchesFromDb,
	}

	if err := client.NewRef(players.Players[playerFromMatch.ID].Name).Set(context.Background(), player); err != nil {
		log.Fatal(err)
	}
}

func GetMatchStats(numberOfMatches string, requesterId string) []database.Match {
	var playerDetails players.PlayerDetails
	_ = mapstructure.Decode(players.Players[requesterId], &playerDetails)

	var player database.Player

	ref := client.NewRef(playerDetails.Name)

	//log.Println("Database key " + ref.Key)
	//log.Println("Database path: " + ref.Path)

	_ = ref.Get(context.Background(), &player)

	var totalMatches, _ = strconv.Atoi(numberOfMatches)
	if totalMatches > len(player.Matches) {
		totalMatches = len(player.Matches)
	}

	for i := 0; i < totalMatches; i++ {
		log.Println("Match " + player.Matches[i].ID + " - Map " + player.Matches[i].Map)
	}

	return player.Matches[:totalMatches]
}

func DeDupeMatches() {
	for player := range players.Players {
		var deDupedMatches = make(map[string]database.Match)
		var deDupedMatchesList = make([]database.Match, 0)
		matchesFromDb := GetMatchStats("10000", player)

		var playerDetails players.PlayerDetails
		_ = mapstructure.Decode(players.Players[player], &playerDetails)

		log.Println("Non sorted matches for " + playerDetails.Name)

		for _, match := range matchesFromDb {
			if _, present := deDupedMatches[match.ID]; !present {
				deDupedMatches[match.ID] = match
				deDupedMatchesList = append(deDupedMatchesList, match)
			}

			log.Println("Match " + match.ID)
		}

		log.Println("De-duped, sorted matches for " + playerDetails.Name)

		for _, match := range deDupedMatchesList {
			log.Println("Match " + match.ID)
		}

		player := database.Player {
			Matches: deDupedMatchesList,
		}

		playerCopy := database.Player {
			Matches: matchesFromDb,
		}

		if err := client.NewRef(playerDetails.Name + "DeDuped").Set(context.Background(), player); err != nil {
			log.Fatal(err)
		}

		if err := client.NewRef(playerDetails.Name + "Copy").Set(context.Background(), playerCopy); err != nil {
			log.Fatal(err)
		}
	}
}