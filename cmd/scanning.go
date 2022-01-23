package cmd

import (
	"context"

	"github.com/dungnh3/guardrails-assignment/internal/apps/job"
	"github.com/dungnh3/guardrails-assignment/pkg/migration"
	"github.com/spf13/cobra"
)

func scanningCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "scanning",
		Short: "start scanning engine",
		Run: func(*cobra.Command, []string) {
			db := cfg.PostgreSQL.ConnectDatabase()

			sqlDB, err := db.DB()
			if err != nil {
				panic(err)
			}
			if err = migration.Up(sqlDB, migrationPath); err != nil {
				logger.Error(err, "failed to do migration")
				panic(err)
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// need this chan to catch error in another routine
			errChan := make(chan error)
			scanningEngine := job.NewScanning(cfg, db)
			go func() {
				if err = scanningEngine.Run(ctx); err != nil {
					logger.Error(err, "scanning job failed")
					errChan <- err
				}
			}()

			handleExitSignal(errChan)
			if err = scanningEngine.Close(ctx); err != nil {
				logger.Error(err, "exception error when shutting down scanning engine")
			}
		},
	}
}
