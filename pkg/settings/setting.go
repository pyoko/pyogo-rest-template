package settings

import (
	"log"

	"github.com/spf13/viper"
)

func ReadConfig(key string) string {
	viper.SetConfigFile("../.env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("Invalid type")
	}

	return value
}