package bootstrap

import (
	"github.com/server-selfish/backend/cmd/di"
	"github.com/server-selfish/backend/config/env"
	"go.uber.org/dig"
)

func Run() *dig.Container {
	env.Load()
	container := di.BuildContainer()
	return container
}
