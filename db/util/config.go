package util

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	DBSourceTest  string `mapstructure:"DB_SOURCE_TEST"`
}

// LoadConfig - Loads config values from config file or env vars
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app") // name of config file (without extension)
	viper.SetConfigType("env") // REQUIRED if the config file does not have the extension in the name

	viper.AutomaticEnv()

	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		err = fmt.Errorf("fatal error reading config: %w", err)
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		err = fmt.Errorf("fatal error unmarshaling config: %w", err)
		return
	}
	return
}
