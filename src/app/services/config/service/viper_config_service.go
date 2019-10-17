package service

import (
	"fmt"
	"services/apps/testing/models/settings"
	"services/apps/testing/services/config"
	"services/package/shared_models"
	"time"

	"github.com/spf13/viper"
)

// Constants for redis cache queries
const (
	CONFIGFILENAME   = "config"
	CONFIGFILETYPE   = "yaml"
	DEFAULTLOCATION1 = "/etc/checkmate/"
	DEFAULTLOCATION2 = "$HOME/.checkmate"
	DEFAULTLOCATION3 = "."
)

type viperConfigService struct {
	application  settings.Application
	cache        settings.Cache
	userService  shared_models.ServiceConnection
	database     shared_models.Database
	cacheTimeout time.Duration
}

// NewViperConfigService will create new an viperConfigService object representation of config.Service interface
func NewViperConfigService() config.Service {

	viper.SetConfigName(CONFIGFILENAME)
	viper.SetConfigType(CONFIGFILETYPE)
	// paths to look for the config file in
	viper.AddConfigPath(DEFAULTLOCATION1)
	viper.AddConfigPath(DEFAULTLOCATION2)
	viper.AddConfigPath(DEFAULTLOCATION3)
	return &viperConfigService{
		application: settings.Application{
			Security:    settings.Security{},
			Concurrency: settings.Concurrency{},
			SMTP:        settings.SMTP{},
		},
		cache: settings.Cache{
			Assessment: shared_models.CacheSettings{},
			User:       shared_models.CacheSettings{},
		},
		userService: shared_models.ServiceConnection{},
		database:    shared_models.Database{},
	}
}

func (vcs *viperConfigService) parseConfig() error {
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	vcs.application.Port = viper.GetString("application.PORT")

	vcs.application.Security.Chiper = viper.GetString("application.security.CHIPER")
	vcs.application.Security.TokenSecret = viper.GetString("application.security.TOKEN_SECRET")
	vcs.application.Concurrency.PoolSize = viper.GetInt("application.concurrency.POOL_SIZE")
	vcs.application.SMTP.User = viper.GetString("application.smtp.USER")
	vcs.application.SMTP.Password = viper.GetString("application.smtp.PASSWORD")

	vcs.userService.Server = viper.GetString("user_service.SERVER")

	vcs.cache.Assessment.Server = viper.GetString("cache.assessment.SERVER")
	vcs.cache.Assessment.Timeout = viper.GetInt("cache.assessment.TIMEOUT_IN_MINUTES")
	vcs.cache.Assessment.MaxActiveConnections = viper.GetInt("cache.assessment.MAX_IDLE_CONNECTIONS")
	vcs.cache.Assessment.MaxIdleConnections = viper.GetInt("cache.assessment.MAX_ACTIVE_CONNECTIONS")

	vcs.cache.User.Server = viper.GetString("cache.user.SERVER")
	vcs.cache.User.Timeout = viper.GetInt("cache.user.TIMEOUT_IN_MINUTES")
	vcs.cache.User.MaxActiveConnections = viper.GetInt("cache.user.MAX_IDLE_CONNECTIONS")
	vcs.cache.User.MaxIdleConnections = viper.GetInt("cache.user.MAX_ACTIVE_CONNECTIONS")

	vcs.database.Host = viper.GetString("database.HOST")
	vcs.database.Port = viper.GetInt("database.PORT")
	vcs.database.Schema = viper.GetString("database.SCHEMA")
	vcs.database.Username = viper.GetString("database.USERNAME")
	vcs.database.Password = viper.GetString("database.PASSWORD")

	return err
}

func (vcs *viperConfigService) Parse() error {
	return vcs.parseConfig()
}

func (vcs *viperConfigService) GetApplicationSettings() settings.Application {
	return vcs.application
}

func (vcs *viperConfigService) GetCacheSettings() settings.Cache {
	return vcs.cache
}

func (vcs *viperConfigService) GetUserServiceSettings() shared_models.ServiceConnection {
	return vcs.userService
}

func (vcs *viperConfigService) GetDatabaseSettings() shared_models.Database {
	return vcs.database
}
