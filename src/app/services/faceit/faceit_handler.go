package faceit

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo"
	"ichor-stats/src/app/models/faceit"
	"ichor-stats/src/app/models/players"
	"ichor-stats/src/app/services/config"
	"ichor-stats/src/app/services/discord/helpers"
	"ichor-stats/src/app/services/firebase"
	"ichor-stats/src/package/api"
	client "ichor-stats/src/package/http"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type ResponseError struct {
	Message string `json:"message"`
}

type FaceitHandler struct {
	FaceitService ServiceFaceit
	ConcurrentLock sync.Mutex
}

type Message struct {
	Body string `json:"body"`
}

func NewFaceitHandler(e *echo.Echo, fs ServiceFaceit) {
	var lockProcessing sync.Mutex

	handler := &FaceitHandler{
		FaceitService: fs,
		ConcurrentLock: lockProcessing,
	}

	g := e.Group("/api/v1/faceit")
	g.POST("/match-end", handler.MatchEnd)
	g.POST("/match-created", handler.MatchCreated)
	g.POST("/match-ready", handler.MatchReady)
	g.POST("/match-configuring", handler.MatchConfiguring)

	c := e.Group("/message")
	c.POST("/custom", handler.CustomMessage)
}

func (fh *FaceitHandler) MatchEnd(c echo.Context) error {
	lock(fh.ConcurrentLock)

	var webhookData = DecipherWebhookData(c)

	req, err := http.NewRequest("GET", api.GetFaceitMatch(webhookData.Payload.MatchID), nil)
	req.Header.Add("Authorization", "Bearer "+ fh.FaceitService.Config.FACEIT_API_KEY)
	response, err := client.Fire(req)
	body, err := ioutil.ReadAll(response.Body)

	log.Println("Match End")
	log.Println(string(body))

	var stats faceit.Match
	_ = json.Unmarshal(body, &stats)

	if err != nil {
		log.Println("Issue decoding finished match - ", err)
	}

	for _, round := range stats.Rounds {
		for _, team := range round.Teams {
			for _, player := range team.Players {
				if playerDetails, playerPresentInMap := players.Players[player.ID]; playerPresentInMap {
					messages := make([]*helpers.Embed, 0)

					matchesFromDb := firebase.GetMatchStats("3", player.ID)

					var uniqueMatch = true

					for _, match := range matchesFromDb {
						if match.ID == webhookData.Payload.MatchID {
							uniqueMatch = false
							log.Println("This match already exists in the database for " + playerDetails.Name + "[" + match.ID + "]")
							break
						}
					}

					if uniqueMatch {
						log.Println("This match is unique, saving to database for " + playerDetails.Name + "[" + webhookData.Payload.MatchID + "]")
						firebase.SaveMatch(player, round, stats, webhookData.Payload.MatchID)
						fh.FaceitService.MatchEnd(player, &messages, stats)
						OutputMessages(fh, &messages)
					}
				}
			}
		}
	}

	time.Sleep(2 * time.Second)
	unlock(fh.ConcurrentLock)

	return c.JSON(http.StatusOK, "")
}

func (fh *FaceitHandler) MatchCreated(c echo.Context) error {
	var messages = make([]*helpers.Embed, 0)
	fh.FaceitService.MatchCreated(DecipherWebhookData(c), &messages)
	OutputMessages(fh, &messages)
	return c.JSON(http.StatusOK, "")
}

func (fh *FaceitHandler) MatchReady(c echo.Context) error {
	var messages = make([]*helpers.Embed, 0)
	fh.FaceitService.MatchReady(DecipherWebhookData(c), &messages)
	OutputMessages(fh, &messages)
	return c.JSON(http.StatusOK, "")
}

func (fh *FaceitHandler) MatchConfiguring(c echo.Context) error {
	var messages = make([]*helpers.Embed, 0)
	fh.FaceitService.MatchConfiguring(DecipherWebhookData(c), &messages)
	OutputMessages(fh, &messages)
	return c.JSON(http.StatusOK, "")
}

func (fh *FaceitHandler) CustomMessage(c echo.Context) error {
	lock(fh.ConcurrentLock)

	var messages = make([]*helpers.Embed, 0)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	body, err := ioutil.ReadAll(c.Request().Body)

	var message Message
	err = json.Unmarshal(body, &message)

	messages = append(messages, helpers.NewEmbed().SetTitle(message.Body))

	if err != nil {
		log.Println(err)
	}

	OutputMessages(fh, &messages)

	time.Sleep(2 * time.Second)
	unlock(fh.ConcurrentLock)

	return c.JSON(http.StatusOK, "")
}

func DecipherWebhookData(c echo.Context) (webhookData faceit.Webhook) {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	body, err := ioutil.ReadAll(c.Request().Body)

	log.Println("Raw body data: " + string(body))

	var webhook faceit.Webhook
	err = json.Unmarshal(body, &webhook)

	if err != nil {
		log.Println(err)
	}

	return webhook
}

func OutputMessages(fh *FaceitHandler, messages *[]*helpers.Embed) {
	if len(*messages) > 0 {
		for _, message := range *messages {
			_, _ = fh.FaceitService.DiscordService.Discord.ChannelMessageSendEmbed(config.GetConfig().CHANNEL_ID, message.MessageEmbed)
		}
	}
}

func lock(mutex sync.Mutex) {
	log.Println("Locking Processing")
}

func unlock(mutex sync.Mutex) {
	log.Println("Unlocking Processing")
}