package mq

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func NewNatsConnection(logger zerolog.Logger) *nats.Conn {
	addr := fmt.Sprintf("%s://%s:%s", viper.GetString("nats.protocol"), viper.GetString("nats.address"), viper.GetString("nats.port"))
	nc, err := nats.Connect(addr,
		nats.UserInfo(viper.GetString("nats.credential.user"), viper.GetString("nats.credential.password")),
		nats.Name(viper.GetString("nats.connetion_name")),
		nats.Timeout(viper.GetDuration("nats.timeout")*time.Second),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(viper.GetDuration("nats.timeout")*time.Second),
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("nats server ping failed")
	}
	return nc
}
