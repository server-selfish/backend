package main

import (
	"github.com/server-selfish/backend/cmd/bootstrap"
	"github.com/server-selfish/backend/cmd/server"
	"github.com/spf13/viper"
)

func main() {
	container := bootstrap.Run()
	httpServer := &server.Server{
		Container: container,
		Address:   ":" + viper.GetString("app.http.port"),
	}
	httpServer.Run()
}
