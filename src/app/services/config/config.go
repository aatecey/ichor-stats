package config

import (
	"fmt"
	"github.com/spf13/viper"
	"ichor-stats/src/app/models/config"
)

const CONFIG_PATH = "./src/build"

func GetConfig() *config.Configuration {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(CONFIG_PATH)
	viper.AutomaticEnv()

	var configuration config.Configuration
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
		panic(1)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
		panic(1)
	}

	return &configuration
}