package faceit

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo"
	"ichor-stats/src/app/models/faceit"
	"log"
	"net/http"
)

type ResponseError struct {
	Message string `json:"message"`
}

type FaceitHandler struct {
	FaceitService ServiceFaceit
}

func NewFaceitHandler(e *echo.Echo, fs ServiceFaceit) {
	handler := &FaceitHandler{
		FaceitService: fs,
	}

	g := e.Group("/api/v1/faceit")
	g.POST("/match-ready", handler.MatchReady)
	g.POST("/match-start", handler.MatchStart)
	g.POST("/match-end", handler.MatchEnd)
}

func (fh *FaceitHandler) MatchEnd(c echo.Context) error {

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var webhook faceit.Webhook
	err := json.NewDecoder(c.Request().Body).Decode(&webhook)
	if err != nil {
		log.Println(err)
		return err
	}

	err = fh.FaceitService.MatchEnd(webhook)
	if err != nil {
		log.Println(err)
	}

	return c.JSON(http.StatusOK, "")
}

func (fh *FaceitHandler) MatchStart(c echo.Context) error {

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var webhook faceit.Webhook
	err := json.NewDecoder(c.Request().Body).Decode(&webhook)
	if err != nil {
		log.Println(err)
		return err
	}

	err = fh.FaceitService.MatchStart(webhook)
	if err != nil {
		log.Println(err)
	}

	return c.JSON(http.StatusOK, "")
}

func (fh *FaceitHandler) MatchReady(c echo.Context) error {

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var webhook faceit.Webhook
	err := json.NewDecoder(c.Request().Body).Decode(&webhook)
	if err != nil {
		log.Println(err)
		return err
	}

	err = fh.FaceitService.MatchReady(webhook)
	if err != nil {
		log.Println(err)
	}

	return c.JSON(http.StatusOK, "")
}
