package main

import (
	"github.com/server-selfish/backend/cmd/bootstrap"
	"github.com/server-selfish/backend/cmd/server"
	"github.com/spf13/viper"
)

// @title           Selfish API
// @version         1.0
// @description     This is selfish api documentation.
// @BasePath        /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Authentication using "Authorization: Bearer <access_token>". The API also accepts the access token from the "selfish_access_token" cookie.
func main() {
	container := bootstrap.Run()
	httpServer := &server.Server{
		Container: container,
		Address:   ":" + viper.GetString("app.http.port"),
	}
	httpServer.Run()
}
