package config

import (
	"log"

	"github.com/spf13/viper"
)

// load config from .env with viper
func Init() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}
