package di

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/server-selfish/backend/config/cache"
	"github.com/server-selfish/backend/config/context"
	docker_client "github.com/server-selfish/backend/config/docker"
	"github.com/server-selfish/backend/config/logger"
	"github.com/server-selfish/backend/config/monitoring"
	"github.com/server-selfish/backend/config/mq"
	"github.com/server-selfish/backend/config/router"
	"github.com/server-selfish/backend/config/storage"
	"github.com/server-selfish/backend/internal/domain/handler"
	container_repository "github.com/server-selfish/backend/internal/domain/repository/container"
	deployment_repository "github.com/server-selfish/backend/internal/domain/repository/deployment"
	github_app_repository "github.com/server-selfish/backend/internal/domain/repository/github_app"
	monitoring_repository "github.com/server-selfish/backend/internal/domain/repository/monitoring"
	project_repository "github.com/server-selfish/backend/internal/domain/repository/project"
	user_repository "github.com/server-selfish/backend/internal/domain/repository/user"
	"github.com/server-selfish/backend/internal/domain/service"
	cache_infra "github.com/server-selfish/backend/internal/infra/cache"
	github_infra "github.com/server-selfish/backend/internal/infra/github"
	monitoring_infra "github.com/server-selfish/backend/internal/infra/monitoring"
	mq_infra "github.com/server-selfish/backend/internal/infra/mq"
	storage_infra "github.com/server-selfish/backend/internal/infra/storage"
	"github.com/server-selfish/backend/internal/pkg"
	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	// logger
	if err := container.Provide(logger.NewLogger); err != nil {
		panic("Failed to provide logger: " + err.Error())
	}

	// app context
	if err := container.Provide(context.NewAppContext); err != nil {
		panic("Failed to provide app context: " + err.Error())
	}
	// token manager
	if err := container.Provide(pkg.NewTokenManager); err != nil {
		panic("Failed to provide token manager: " + err.Error())
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

	// prometheus connection
	if err := container.Provide(monitoring.NewMonitoring); err != nil {
		panic("Failed to provide monitoring instance: " + err.Error())
	}

	// db connection
	if err := container.Provide(storage.NewPostgresqlConn); err != nil {
		panic("Failed to provide db connection: " + err.Error())
	}
	if err := container.Provide(func(pool *pgxpool.Pool) project_repository.DBTX { return pool }); err != nil {
		panic("Failed to provide project dbtx: " + err.Error())
	}
	if err := container.Provide(func(pool *pgxpool.Pool) deployment_repository.DBTX { return pool }); err != nil {
		panic("Failed to provide deployment dbtx: " + err.Error())
	}
	if err := container.Provide(func(pool *pgxpool.Pool) container_repository.DBTX { return pool }); err != nil {
		panic("Failed to provide container dbtx: " + err.Error())
	}
	if err := container.Provide(func(pool *pgxpool.Pool) user_repository.DBTX { return pool }); err != nil {
		panic("Failed to provide user dbtx: " + err.Error())
	}
	if err := container.Provide(func(pool *pgxpool.Pool) github_app_repository.DBTX { return pool }); err != nil {
		panic("Failed to provide github app dbtx: " + err.Error())
	}

	// valkey connection
	if err := container.Provide(cache.NewValkeyConnection); err != nil {
		panic("Failed to provide cache connection: " + err.Error())
	}

	if err := container.Provide(docker_client.NewDockerClient); err != nil {
		panic("Failed to provide docker client: " + err.Error())
	}

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
	if err := container.Provide(github_infra.NewGithubInfra); err != nil {
		panic("Failed to provide github infra: " + err.Error())
	}
	if err := container.Provide(monitoring_infra.NewPrometheusInfra); err != nil {
		panic("Failed to provide prometheus infra: " + err.Error())
	}

	// utils
	if err := container.Provide(pkg.NewTxManager); err != nil {
		panic("Failed to provide tx manager: " + err.Error())
	}

	// repositories
	if err := container.Provide(project_repository.New); err != nil {
		panic("Failed to provide project repository: " + err.Error())
	}
	if err := container.Provide(deployment_repository.New); err != nil {
		panic("Failed to provide deployment repository: " + err.Error())
	}
	if err := container.Provide(container_repository.New); err != nil {
		panic("Failed to provide container repository: " + err.Error())
	}
	if err := container.Provide(container_repository.NewContainerRepository); err != nil {
		panic("Failed to provide container log repository : " + err.Error())
	}
	if err := container.Provide(user_repository.New); err != nil {
		panic("Failed to provide user repository: " + err.Error())
	}
	if err := container.Provide(github_app_repository.New); err != nil {
		panic("Failed to provide github app repository: " + err.Error())
	}
	if err := container.Provide(monitoring_repository.NewPrometheusRepository); err != nil {
		panic("Failed to provide monitoring repository: " + err.Error())
	}

	// services
	if err := container.Provide(service.NewProjectService); err != nil {
		panic("Failed to provide Project Service: " + err.Error())
	}
	if err := container.Provide(service.NewAuthService); err != nil {
		panic("Failed to provide Auth Service: " + err.Error())
	}
	if err := container.Provide(service.NewContainerService); err != nil {
		panic("Failed to provide Container Service: " + err.Error())
	}
	if err := container.Provide(service.NewGithubAppService); err != nil {
		panic("Failed to provide Github App Service: " + err.Error())
	}
	if err := container.Provide(service.NewDeploymentService); err != nil {
		panic("Failed to provide Deployment Service: " + err.Error())
	}
	if err := container.Provide(service.NewPrometheusService); err != nil {
		panic("Failed to provide prometheus Service: " + err.Error())
	}

	// handlers
	if err := container.Provide(handler.NewProjectHandler); err != nil {
		panic("Failed to provide Project Handler: " + err.Error())
	}
	if err := container.Provide(handler.NewDeploymentHandler); err != nil {
		panic("Failed to provide Deployment Handler: " + err.Error())
	}
	if err := container.Provide(handler.NewAuthHandler); err != nil {
		panic("Failed to provide Auth Handler: " + err.Error())
	}
	if err := container.Provide(handler.NewGithubAppHandler); err != nil {
		panic("Failed to provide Github App Handler: " + err.Error())
	}
	if err := container.Provide(handler.NewContainerHandler); err != nil {
		panic("Failed to provide Container Handler: " + err.Error())
	}
	if err := container.Provide(handler.NewMonitoringHandler); err != nil {
		panic("Failed to provide monitoring Handler: " + err.Error())
	}

	// http server
	if err := container.Provide(router.NewHTTPChi); err != nil {
		panic("Failed to provide http server: " + err.Error())
	}

	return container
}
