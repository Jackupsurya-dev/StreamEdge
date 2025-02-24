package utils

import (
	"producer/constants"
	"producer/logger"

	"github.com/spf13/viper"
)

// Read Configuration files
func Configuration(configPath string) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName(constants.CONFIG_FILE_NAME)

	err := viper.ReadInConfig()
	if err != nil {
		logger.Log.Errorln("Failed to Start Web Server: Viper Read Config Failure -", err)
		panic(err)
	}

	viper.WatchConfig()
}
