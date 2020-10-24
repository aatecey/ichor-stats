package faceit

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo"
	"ichor-stats/src/app/models/faceit"
	"ichor-stats/src/app/services/config"
	"ichor-stats/src/app/services/discord/helpers"
	"io/ioutil"
	"log"
	"net/http"
)

type ResponseError struct {
	Message string `json:"message"`
}

type FaceitHandler struct {
	FaceitService ServiceFaceit
}

type Message struct {
	Body string `json:"body"`
}

func NewFaceitHandler(e *echo.Echo, fs ServiceFaceit) {
	handler := &FaceitHandler{
		FaceitService: fs,
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
	var messages = make([]*helpers.Embed, 0)
	fh.FaceitService.MatchEnd(DecipherWebhookData(c), &messages)
	OutputMessages(fh, &messages)
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