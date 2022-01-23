package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dungnh3/guardrails-assignment/config"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
)

var (
	cfg    *config.Config
	logger logr.Logger
)

const migrationPath = "file://db/migrations"

func Run(args []string) error {
	var err error
	cfg, err = config.Load()
	if err != nil {
		return err
	}
	logger = cfg.Logger.MustBuildLogR()

	command := &cobra.Command{
		Use:          "guardrails",
		Short:        "guardrails platform",
		SilenceUsage: true,
	}
	command.AddCommand(serverCommand())
	command.AddCommand(scanningCmd())
	return command.Execute()
}

func handleExitSignal(errChan <-chan error) {
	quit := catchInterruptSignal()
	for {
		select {
		case <-quit:
			return
		case <-errChan:
			return
		}
	}
}

func catchInterruptSignal() <-chan os.Signal {
	// wait for interrupt signal here
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	return quit
}
