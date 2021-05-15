package cmd

import (
	"github.com/spf13/cobra"
)

var migrateCMD = cobra.Command{
	Use:     "migrate",
	Long:    "migrate database strucutures. This will migrate tables",
	Aliases: []string{"m"},
	Run:     Runner.migrate,
}

// migrate database with fake data
func (c *command) migrate(cmd *cobra.Command, args []string) {}
