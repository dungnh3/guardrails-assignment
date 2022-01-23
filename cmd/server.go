package cmd

import (
	server "github.com/dungnh3/guardrails-assignment/internal/apps/server"
	"github.com/dungnh3/guardrails-assignment/pkg/migration"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func serverCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "start ims-olympus server",
		RunE: func(cmd *cobra.Command, args []string) error {
			zapLogger, err := cfg.Logger.Build()
			if err != nil {
				return err
			}

			grpcMiddlewareUnary := []grpc.UnaryServerInterceptor{
				grpc_prometheus.UnaryServerInterceptor,
				grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
				grpc_zap.UnaryServerInterceptor(zapLogger),
				grpc_validator.UnaryServerInterceptor(),
			}

			s := server.NewServer(
				cfg.Server,
				grpc_middleware.WithUnaryServerChain(grpcMiddlewareUnary...),
			)

			db := cfg.PostgreSQL.ConnectDatabase()

			sqlDB, err := db.DB()
			if err != nil {
				panic(err)
			}
			if err = migration.Up(sqlDB, migrationPath); err != nil {
				logger.Error(err, "failed to do migration")
				panic(err)
			}

			svc, err := server.NewService(cfg, db)
			if err != nil {
				return err
			}
			defer func() {
				if err = svc.Close(); err != nil {
					logger.Error(err, "exception error when closing service", "error", err.Error())
					return
				}
				logger.Info("closing service success")
			}()

			if err = s.Register(svc); err != nil {
				logger.Error(err, "error register servers")
				return err
			}

			logger.Info("Starting server", "grpc", cfg.Server.GRPC, "http", cfg.Server.HTTP)
			if err = s.Serve(); err != nil {
				logger.Error(err, "error start server")
				return err
			}
			return nil
		},
	}
}
