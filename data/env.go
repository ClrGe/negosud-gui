package data

import (
	"github.com/spf13/viper"
)

// --------------------------- ENVIRONMENT ------------------------------

// define and load env. variables contained in app.env

type Config struct {
	ENV         string `mapstructure:"ENV"`
	SERVER_DEV  string `mapstructure:"SERVER_DEV"`
	SERVER_PROD string `mapstructure:"SERVER_PROD"`
	KEY         string `mapstructure:"KEY"`
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

func SaveConfig(key string, value string) {
	viper.Set(key, value)
	viper.WriteConfig()

	return

}
