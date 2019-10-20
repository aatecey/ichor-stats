package api

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo"
	"ichor-stats/src/app/services/discord"
	"log"
	"net/http"
)

type ResponseError struct {
	Message string `json:"message"`
}

func NewFaceitHandler(e *echo.Echo) {
	g := e.Group("/api/v1/faceit")
	g.POST("/match-end", MatchEnd)
}

const (
	FACEIT_API_KEY = "***"
	DISCORD_BOT_ID = "***"
	CHANNEL_ID     = "***"

)

func MatchEnd(c echo.Context) error {

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var webhook Webhook
	err := json.NewDecoder(c.Request().Body).Decode(&webhook)
	if err != nil {
		log.Println(err)
		return err
	}

	url := "https://open.faceit.com/data/v4/matches/" + webhook.Payload.MatchID + "/stats"

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + FACEIT_API_KEY

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	var stats Match
	err = json.NewDecoder(resp.Body).Decode(&stats)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(&stats)
	log.Println(stats)
	for _, s := range stats.Rounds {
		for _, a := range s.Teams {
			log.Println(a.MatchStats.Map)
			log.Println(a.MatchStats.Score)
			for _, d := range a.Players {
				if d.ID == "0d94613d-b736-46ba-b8cd-d2159ddad705" || d.ID == "b26df7d4-8517-4ec6-ab58-708487e5fe60" || d.ID == "b0a57a5a-2f7a-481c-aaa8-8013a83378e3" {
					log.Println(d.Stats.Kills)
					log.Println(d.Stats.Assists)
					log.Println(d.Stats.Deaths)
					log.Println(d.Stats.KD)
					log.Println(d.Stats.KR)
				}
			}
		}
	}

	discord.SendMessage("Match ended")

	return c.JSON(http.StatusOK, "")
}

type Webhook struct {
	Payload Payload `json:"payload"`
}

type Payload struct {
	MatchID string `json:"id"`
}

type Match struct {
	Rounds []Rounds `json:"rounds"`
}

type Rounds struct {
	MatchStats MatchStats `json:"round_stats"`
	Teams []Teams `json:"teams"`
}

type MatchStats struct {
	Map string `json:"Map"`
	Score string `json:"Score"`
}

type Teams struct {
	MatchStats MatchStats `json:"round_stats"`
	Players []Players `json:"players"`
}

type Players struct {
	ID string `json:"player_id"`
	Stats PlayerStats `json:"player_stats"`
}

type PlayerStats struct {
	Kills string `json:"Kills"`
	Assists string `json:"Assists"`
	Deaths string `json:"Deaths"`
	KD string `json:"K/D Ratio"`
	KR string `json:"K/R Ratio"`
}
