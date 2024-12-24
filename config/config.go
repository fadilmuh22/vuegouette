package config

import (
	"github.com/spf13/viper"
)

// load config from .env with viper
func Init() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
}
