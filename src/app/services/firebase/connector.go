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

func SaveMatch(matchStats faceit.Match, matchId string) {
	for _, s := range matchStats.Rounds {
		for _, a := range s.Teams {
			for _, d := range a.Players {
				if playerDetails, playerPresentInMap := players.Players[d.ID]; playerPresentInMap {

					matchesFromDb := GetMatchStats("10000", d.ID)

					var outcome = "Win"

					if d.Stats.Result == "0" {
						outcome = "Loss"
					}

					match := database.Match{
						ID:                 matchId,
						Map:                s.MatchStats.Map,
						Result:             outcome,
						Score:              matchStats.Rounds[0].MatchStats.Score,
						Kills:              d.Stats.Kills,
						Assists:            d.Stats.Assists,
						Deaths:             d.Stats.Deaths,
						Headshots:          d.Stats.Headshots,
						HeadshotPercentage: d.Stats.HeadshotPercentage,
						Pentas:             d.Stats.Pentas,
						Quads:              d.Stats.Quads,
						Triples:            d.Stats.Triples,
						KillDeathRatio:     d.Stats.KD,
						KillRoundRatio:     d.Stats.KR,
						MVPs:               d.Stats.MVPs,
					}

					matchesFromDb = append([]database.Match{match}, matchesFromDb...)

					player := database.Player{
						Matches: matchesFromDb,
					}

					if err := client.NewRef(playerDetails.Name).Set(context.Background(), player); err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}
}

func GetMatchStats(numberOfMatches string, requesterId string) []database.Match {
	var playerDetails players.PlayerDetails
	_ = mapstructure.Decode(players.Players[requesterId], &playerDetails)

	var player database.Player

	ref := client.NewRef(playerDetails.Name)

	log.Println("Database key " + ref.Key)
	log.Println("Database path: " + ref.Path)

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