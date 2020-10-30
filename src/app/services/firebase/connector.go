package firebase

import (
	"context"
	"encoding/json"
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
			matchErr := json.Unmarshal(api.FaceitRequest(api.GetFaceitMatchDetails(match.MatchId)), &matchDetails)
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

func SaveMatch(playerFromMatch faceit.Players, round faceit.Rounds, matchStats faceit.Match, matchId string) {
	matchesFromDb := GetMatchStats("10000", playerFromMatch.ID)

	var outcome = "Win"

	if playerFromMatch.Stats.Result == "0" {
		outcome = "Loss"
	}

	match := database.Match{
		ID:                 matchId,
		Map:                round.MatchStats.Map,
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

	_ = ref.Get(context.Background(), &player)

	var totalMatches, _ = strconv.Atoi(numberOfMatches)
	if totalMatches > len(player.Matches) {
		totalMatches = len(player.Matches)
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

func RetrospectiveUpdate() {
	for player := range players.Players {

		matchesViaFaceItApi := make([]database.Match, 0)
		matchesFromDb := GetMatchStats("10000", player)
		combinedMatches := make([]database.Match, 0)

		var playerDetails players.PlayerDetails
		_ = mapstructure.Decode(players.Players[player], &playerDetails)

		log.Println("Requester = " + player)
		matchHistory, _ := helpers.GetMatchHistory("20", player)

		for _, match := range matchHistory.MatchItem {
			var matchDetails faceit.Match
			matchErr := json.Unmarshal(api.FaceitRequest(api.GetFaceitMatchDetails(match.MatchId)), &matchDetails)
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

			matchesViaFaceItApi = append(matchesViaFaceItApi, match)
			log.Println("Length = " + string(len(matchesViaFaceItApi)))

		}

		log.Println("Length = " + string(len(matchesViaFaceItApi)))

		for _, dbMatch := range matchesFromDb {
			retrySameMatch := true
			for retrySameMatch {
				var faceitApiMatch = database.Match{}

				if len(matchesViaFaceItApi) > 0 {
					faceitApiMatch = matchesViaFaceItApi[0]
				}

				log.Println("FaceIt Match ID : " + faceitApiMatch.ID)
				log.Println("DB Match ID : " + dbMatch.ID)

				if faceitApiMatch != (database.Match{}) && dbMatch.ID != faceitApiMatch.ID {
					combinedMatches = append(combinedMatches, faceitApiMatch)
				} else {
					combinedMatches = append(combinedMatches, dbMatch)
					retrySameMatch = false
				}

				if len(matchesViaFaceItApi) > 0 {
					matchesViaFaceItApi = matchesViaFaceItApi[1:]
				}
			}
		}

		player := database.Player{
			Matches: combinedMatches,
		}

		if err := client.NewRef(playerDetails.Name + "RetroUpdate").Set(context.Background(), player); err != nil {
			log.Fatal(err)
		}
	}
}