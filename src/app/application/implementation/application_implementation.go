package implementation

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/tylerb/graceful"
	"ichor-stats/src/app/application"
	"ichor-stats/src/app/services/config"
	"ichor-stats/src/app/services/discord"
	"ichor-stats/src/app/services/faceit"
	"ichor-stats/src/app/services/firebase"
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
	initializeServices(echo)
	initializeMiddleWare(echo)

	loadedConfig := config.GetConfig()
	log.Println(loadedConfig.DISCORD_BOT_ID)
	log.Println(loadedConfig.CHANNEL_ID)
	log.Println(loadedConfig.FACEIT_API_KEY)

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
	return echo
}

func initializeServices(echo *echo.Echo) {
	appConfig := config.GetConfig()

	firebase.Init()

	//Call to update database with faceit games that didnt save correctly
	//firebase.RetrospectiveUpdate()

	//Call to de-duplicate data in database if it ever gets corrupt.
	//firebase.DeDupeMatches()

	//Call to initialise data in database if it ever gets corrupt.
	//firebase.Setup()

	discordService := discord.NewDiscordService(appConfig)
	discord.NewDiscordHandler(&discordService, appConfig)

	faceitService := faceit.NewFaceitService(appConfig, discordService)
	faceit.NewFaceitHandler(echo, faceitService)
}
