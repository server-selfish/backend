package storage

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func NewRustfsConnection(logger zerolog.Logger) *minio.Client {
	hostPort := viper.GetString("rustfs.host") + ":" + viper.GetString("rustfs.port")
	rustfsClient, err := minio.New(hostPort, &minio.Options{
		Creds: credentials.NewStaticV4(
			viper.GetString("rustfs.credential.user"),
			viper.GetString("rustfs.credential.password"),
			"",
		),
		Secure: false,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("rustfs connection instance error")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rustfsClient.ListBuckets(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("rustfs ping failed")
	}

	return rustfsClient
}
