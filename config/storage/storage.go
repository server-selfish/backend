package storage

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

func NewRustfsConnection() *minio.Client {
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
		log.Fatalln(err)
	}
	return rustfsClient
}
