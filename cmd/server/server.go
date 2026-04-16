package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moby/moby/client"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/server-selfish/backend/internal/domain/handler"
	"github.com/server-selfish/backend/internal/domain/service"
	"github.com/valkey-io/valkey-go"
	"go.uber.org/dig"
)

type Server struct {
	Container   *dig.Container
	ServerReady chan bool
	Address     string
}

func (s *Server) Run(ctx context.Context) {
	err := s.Container.Invoke(
		func(
			logger zerolog.Logger,
			r chi.Router,
			cache valkey.Client,
			mq *nats.Conn,
			db *pgxpool.Pool,
			ph handler.ProjectHandler,
			dh handler.DeploymentHandler,
			ah handler.AuthHandler,
			as service.AuthService,
			dc *client.Client,
		) {
			defer cache.Close()
			defer dc.Close()

			defer func() {
				if err := mq.Drain(); err != nil {
					logger.Error().Err(err).Msg("Failed to close mq client connection")
				}
			}()
			defer db.Close()

			// Public auth routes
			handler.RegisterPublicAuthRoutes(r, ah)

			// Protected business routes
			r.Group(func(pr chi.Router) {
				pr.Use(handler.RequireAuth(as))

				handler.RegisterProtectedAuthRoutes(pr, ah)
				handler.RegisterProjectRoutes(pr, ph)
				handler.RegisterDeploymentRoutes(pr, dh)
			})

			srv := &http.Server{
				Addr:              s.Address,
				Handler:           r,
				ReadHeaderTimeout: 5 * time.Second,
			}

			go func() {
				e := os.Getenv("ENV")
				switch e {
				case "production":
					if err := srv.ListenAndServeTLS("./config/cert/server.crt", "./config/cert/server.key"); err != nil && err != http.ErrServerClosed {
						logger.Fatal().Err(err).Msg("Failed to listen and serve http server")
					}
				default:
					if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						logger.Fatal().Err(err).Msg("Failed to listen and serve http server")
					}
				}
			}()

			if s.ServerReady != nil {
				for range 50 {
					conn, err := net.DialTimeout("tcp", s.Address, 100*time.Millisecond)
					if err == nil {
						if err := conn.Close(); err != nil {
							logger.Fatal().Err(err).Msg("establish check connection failed to close")
						}
						s.ServerReady <- true
						break
					}
					time.Sleep(100 * time.Millisecond)
				}
			}

			logger.Info().Msgf("HTTP Server Starting in port %s", s.Address)
			<-ctx.Done()

			logger.Info().Msg("Shutting down server...")

			shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			if err := srv.Shutdown(shutdownCtx); err != nil {
				logger.Fatal().Err(err).Msg("HTTP Server forced to shutdown")
			}

			logger.Info().Msg("Server exiting...")
		},
	)
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}
}
