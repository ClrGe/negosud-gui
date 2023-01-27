package data

import (
	"github.com/spf13/viper"
)

// --------------------------- ENVIRONMENT ------------------------------

// define and load env. variables contained in app.env

type Config struct {
	SERVER string `mapstructure:"SERVER"`
	APIKEY string `mapstructure:"APIKEY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)

	return
}
