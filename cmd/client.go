package main

import (
	"github.com/DVKunion/SeaMoon/pkg/client"
	"github.com/DVKunion/SeaMoon/pkg/consts"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var (
	debug   bool
	verbose bool

	rootCommand = &cobra.Command{
		Run: Client,
	}

	clientCommand = &cobra.Command{
		Use:   "client",
		Short: "SeaMoon Client",
		Run:   Client,
	}

	versionCommand = &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info(consts.Version)
		},
	}
)

func init() {
	rootCommand.AddCommand(versionCommand)
	clientCommand.Flags().BoolVarP(&verbose, "verbose", "v", false, "proxy detail log")
	clientCommand.Flags().BoolVarP(&debug, "debug", "d", false, "proxy detail log")

	rootCommand.AddCommand(clientCommand)
}

func Client(cmd *cobra.Command, args []string) {
	client.Serve(cmd.Context(), verbose, debug)
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
