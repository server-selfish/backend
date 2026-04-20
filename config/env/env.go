package env

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// config.yaml that loaded by this code in production environment (container)
// can be achieved by mount the config file to the container at / (default) or read the config file
// from remote source, which can you do based on docs (https://github.com/spf13/viper)

// or (which i don't suggest but working) inject the config.yaml in the Dockerfile.

func Load() {
	_ = godotenv.Load()
	env := os.Getenv("ENV")
	configName := "config.local"
	configpath := "."
	switch env {
	case "production":
		configName = "config"
	case "test":
		configName = "config.test"
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configpath)
	if err := viper.ReadInConfig(); err != nil {
		panic("failed to read config")
	}
}
