package implementation

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/tylerb/graceful"
	"ichor-stats/src/app/application"
	"ichor-stats/src/app/services/api"
	"ichor-stats/src/app/services/discord"
	"log"
	"net/http"
	"time"
)

type app struct {
}

// NewApplication will create a new application object representation of package.Application interface
func NewApplication() application.Application {
	return &app{}
}

func (a *app) Run() {
	echo := initialize()
	initializeMiddleWare(echo)
	log.Println("Service listening on " + echo.Server.Addr)

	err := graceful.ListenAndServe(echo.Server, 5 * time.Second)
	if err != nil {
		fmt.Println(err)
	}
}

func initializeMiddleWare(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://ichor-stats.azurewebsites.net"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{"Accept", "Accept-Language", "Content-Type"},
	}))
}

func initialize() *echo.Echo {

	echo := echo.New()
	echo.Server.Addr = ":" + "5000"

	api.NewFaceitHandler(echo)

	discordgo, err := discordgo.New("Bot " + "NjM0NzI2MTI1Nzk2NDU4NTA3.XauI2A.WgtyCNjHY_kvnPjVdLMepNFjeFg")
	if err != nil {
		fmt.Println(err)
	}

	// Register messageCreate as a callback for the messageCreate events.
	discordgo.AddHandler(discord.MessageCreate)

	// Open the websocket and begin listening.
	err = discordgo.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	} else {
		fmt.Println("Discord websocket open")
	}

	return echo
}
