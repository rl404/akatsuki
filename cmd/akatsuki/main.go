package main

import (
	_ "github.com/rl404/akatsuki/docs"
	"github.com/rl404/akatsuki/internal/utils"

	"github.com/spf13/cobra"
)

// @title Akatsuki API
// @description Akatsuki API.
// @BasePath /
// @schemes http https
func main() {
	cmd := cobra.Command{
		Use:   "akatsuki",
		Short: "Akatsuki API",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "Run API server",
		RunE: func(*cobra.Command, []string) error {
			return server()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "migrate",
		Short: "Migrate database",
		RunE: func(*cobra.Command, []string) error {
			return migrate()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "consumer",
		Short: "Run message consumer",
		RunE: func(*cobra.Command, []string) error {
			return consumer()
		},
	})

	cronCmd := cobra.Command{
		Use:   "cron",
		Short: "Cron",
	}

	cronCmd.AddCommand(&cobra.Command{
		Use:   "update",
		Short: "Update old data",
		RunE: func(*cobra.Command, []string) error {
			return cronUpdate()
		},
	})

	cmd.AddCommand(&cronCmd)

	if err := cmd.Execute(); err != nil {
		utils.Fatal(err.Error())
	}
}
