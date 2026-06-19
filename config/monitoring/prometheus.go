package monitoring

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func NewMonitoring(logger zerolog.Logger) (v1.API, error) {
	prometheusURL := viper.GetString("prometheus.url")
	client, err := api.NewClient(api.Config{
		Address: prometheusURL,
		RoundTripper: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("failed create prometheus connection")
		return nil, err
	}

	v1api := v1.NewAPI(client)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, _, err = v1api.Query(ctx, "up", time.Now())
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to ping to prometheus")
		return nil, err
	}
	return v1api, nil
}
