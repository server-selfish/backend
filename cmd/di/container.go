package di

import (
	"github.com/nats-io/nats.go/jetstream"
	"github.com/server-selfish/backend/config/cache"
	"github.com/server-selfish/backend/config/logger"
	"github.com/server-selfish/backend/config/mq"
	"github.com/server-selfish/backend/config/router"
	"github.com/server-selfish/backend/config/storage"
	"github.com/server-selfish/backend/internal/domain/handler"
	project_repository "github.com/server-selfish/backend/internal/domain/repository/project"
	"github.com/server-selfish/backend/internal/domain/service"
	cache_infra "github.com/server-selfish/backend/internal/infra/cache"
	mq_infra "github.com/server-selfish/backend/internal/infra/mq"
	storage_infra "github.com/server-selfish/backend/internal/infra/storage"
	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	// here is your dependency injection using dig with order matter,
	// but for the very first time, i will give you a default one.
	// you can change this anytime

	// logger
	if err := container.Provide(logger.NewLogger); err != nil {
		panic("Failed to provide logger: " + err.Error())
	}
	// object storage connection
	if err := container.Provide(storage.NewRustfsConnection); err != nil {
		panic("Failed to provide object storage connection: " + err.Error())
	}

	// mq client connection
	if err := container.Provide(mq.NewNatsConnection); err != nil {
		panic("Failed to provide mq connection: " + err.Error())
	}
	// jetstream connection
	if err := container.Provide(jetstream.New); err != nil {
		panic("Failed to provide jetstream instance: " + err.Error())
	}

	// db connection
	if err := container.Provide(storage.NewPostgresqlConn); err != nil {
		panic("Failed to provide db connection: " + err.Error())
	}
	//  connection
	if err := container.Provide(cache.NewValkeyConnection); err != nil {
		panic("Failed to provide cache connection: " + err.Error())
	}

	// you can add your own handler, service, repository,infra, or even
	// your own defined config here and invoke in the /cmd/server/http_server.go

	// infra
	if err := container.Provide(cache_infra.NewValkeyCache); err != nil {
		panic("Failed to provide cache infra: " + err.Error())
	}
	if err := container.Provide(mq_infra.NewJetstreamInfra); err != nil {
		panic("Failed to provide MQ infra: " + err.Error())
	}
	if err := container.Provide(storage_infra.NewRustfsInfra); err != nil {
		panic("Failed to provide object storage infra: " + err.Error())
	}

	// repo
	if err := container.Provide(project_repository.New); err != nil {
		panic("Failed to provide Project repository: " + err.Error())
	}

	// service
	if err := container.Provide(service.NewProjectService); err != nil {
		panic("Failed to provide Project Service: " + err.Error())
	}

	// handler
	if err := container.Provide(handler.NewProjectHandler); err != nil {
		panic("Failed to provide Project Handler: " + err.Error())
	}

	// http server
	if err := container.Provide(router.NewHTTPChi); err != nil {
		panic("Failed to provide http server: " + err.Error())
	}
	return container
}
