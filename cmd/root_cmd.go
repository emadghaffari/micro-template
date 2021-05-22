package cmd

import (
	"micro/app"
	"micro/config"

	"github.com/spf13/cobra"
)

var (
	Runner     CommandLine = &command{}
	configFile             = ""
	debug      bool
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
		config.Confs.SetDebug(debug)
		app.Base.StartApplication()
	},
}

// RootCmd will add flags and subcommands to the different commands
func (c *command) RootCmd() *cobra.Command {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "The configuration file")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "The service debug(true is production - false is dev)")

	// add more commands
	rootCmd.AddCommand(&migrateCMD)
	rootCmd.AddCommand(&seedCMD)
	return &rootCmd
}
