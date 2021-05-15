package cmd

import (
	"micro/app"

	"github.com/spf13/cobra"
)

var (
	Runner     CommandLine = &command{}
	configFile             = ""
)

type CommandLine interface {
	RootCmd() *cobra.Command
	migrate(cmd *cobra.Command, args []string)
	seed(cmd *cobra.Command, args []string)
}

type command struct{}

// rootCmd will run the log streamer
var rootCmd = cobra.Command{
	Use:  "micro",
	Long: "A service that will validate restful transactions and send them to stripe.",
	Run: func(cmd *cobra.Command, args []string) {
		app.Base.StartApplication()
	},
}

// RootCmd will add flags and subcommands to the different commands
func (c *command) RootCmd() *cobra.Command {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "The configuration file")

	// add more commands
	rootCmd.AddCommand(&seedCMD)
	rootCmd.AddCommand(&migrateCMD)
	return &rootCmd
}
