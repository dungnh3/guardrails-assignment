package cmd

import (
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
		Use:          "ims-olympus",
		Short:        "ims olympus platform",
		SilenceUsage: true,
	}
	command.AddCommand(serverCommand())
	return command.Execute()
}
