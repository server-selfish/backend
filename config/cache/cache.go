package cache

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/valkey-io/valkey-go"
)

func NewValkeyConnection(logger zerolog.Logger) (valkey.Client, error) {
	addr := fmt.Sprintf("%s:%s", viper.GetString("valkey.address"), viper.GetString("valkey.port"))
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{addr},
		Username:    viper.GetString("valkey.username"),
		Password:    viper.GetString("valkey.password"),
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to valkey")
	}
	return client, nil
}
